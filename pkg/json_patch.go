package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

type JsonPatch struct {
	Op    string          `json:"op"`
	Path  string          `json:"path"`
	From  string          `json:"from,omitempty"`
	Value json.RawMessage `json:"value,omitempty"`
}

type PreparePatchCtx struct {
	Kind            string
	Path            ParsedPath
	Model           any
	Value           any
	SerializedPatch []byte
	Generic         map[string]any
	Content         []byte
}

func parsePathStep(patch JsonPatch) Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Parse patch",
		Run: func(ctx *PreparePatchCtx) error {
			var ierr error
			ctx.Path, ierr = ParsePath(patch.Path)
			return ierr
		},
	}
}

func typeFromEntitiesStep(ctx context.Context, db *bun.DB) Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Extract type from entities",
		Run: func(pctx *PreparePatchCtx) error {
			return db.NewSelect().Model((*models.Entity)(nil)).Where("mrid = ?", pctx.Path.Mrid).Column("entity_type").Scan(ctx, &pctx.Kind)
		},
	}
}

func formTypeStep() Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Extract form type",
		Run: func(ctx *PreparePatchCtx) error {
			kinds := FormTypes()
			var ok bool
			ctx.Model, ok = kinds[ctx.Kind]
			return ErrorIfNotOk(ok, fmt.Sprintf("Unknown form type %s", ctx.Kind))
		},
	}
}

func extractLastEntryStep(ctx context.Context, db *bun.DB) Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Extract last entity",
		Run: func(pctx *PreparePatchCtx) error {
			return db.NewSelect().Model(pctx.Model).Where("mrid = ?", pctx.Path.Mrid).OrderBy("commit_id", bun.OrderDesc).Limit(1).Scan(ctx)
		},
	}
}

func intepretValueStep(patch JsonPatch) Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Interprete value",
		Run: func(ctx *PreparePatchCtx) error {
			return json.Unmarshal(patch.Value, &ctx.Value)
		},
	}
}

func serializePatchStep(patch JsonPatch) Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Serialize patch step",
		Run: func(ctx *PreparePatchCtx) error {
			var ierr error
			ctx.SerializedPatch, ierr = json.Marshal(patch)
			return ierr
		},
	}
}

func applyPatchStep(patch JsonPatch) Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Apply patch",
		Run: func(ctx *PreparePatchCtx) error {
			ctx.Content = Must(json.Marshal(ctx.Model))
			PanicOnErr(json.Unmarshal(ctx.Content, &ctx.Generic))

			switch patch.Op {
			case "replace":
				ctx.Generic[ctx.Path.Field] = ctx.Value
			default:
				return fmt.Errorf("Unsupported operation %s", patch.Op)
			}
			return nil
		},
	}
}

func serializingGenericModelStep() Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Serializing updated generic map",
		Run: func(ctx *PreparePatchCtx) error {
			var ierr error
			ctx.Content, ierr = json.Marshal(ctx.Generic)
			return ierr
		},
	}
}

func updatingOriginalModelStep() Step[PreparePatchCtx] {
	return Step[PreparePatchCtx]{
		Name: "Updating original model",
		Run: func(ctx *PreparePatchCtx) error {
			return json.Unmarshal(ctx.Content, ctx.Model)
		},
	}
}

func ApplyPatch(ctx context.Context, db *bun.DB, patch JsonPatch) error {
	var prepCtx PreparePatchCtx
	err := Pipe(&prepCtx,
		parsePathStep(patch),
		typeFromEntitiesStep(ctx, db),
		formTypeStep(),
		extractLastEntryStep(ctx, db),
		intepretValueStep(patch),
		serializePatchStep(patch),
		applyPatchStep(patch),
		serializingGenericModelStep(),
		updatingOriginalModelStep(),
	)

	if err != nil {
		return fmt.Errorf("Prepare data failed: %w", err)
	}

	commit := models.Commit{
		Message: fmt.Sprintf("Applied json patch %s", string(prepCtx.SerializedPatch)),
		Author:  "Json patcher",
	}

	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := ReturnOnFirstError(
			func() error {
				_, ierr := tx.NewInsert().Model(&commit).Exec(ctx)
				return ierr
			},
			func() error {
				setCommitIfPossible(prepCtx.Model, int(commit.Id))
				_, ierr := tx.NewInsert().Model(prepCtx.Model).ExcludeColumn("id").Exec(ctx)
				return ierr
			},
		)
		return err
	})

}

type ParsedPath struct {
	Mrid  string
	Field string
}

func ParsePath(path string) (ParsedPath, error) {
	result := strings.Split(path, "/")
	if len(result) != 3 {
		return ParsedPath{}, fmt.Errorf("Path %s must match /<resoureId>/<field>", path)
	}
	return ParsedPath{Mrid: result[1], Field: result[2]}, nil
}

func ErrorIfNotOk(ok bool, msg string) error {
	if !ok {
		return fmt.Errorf("Wanted true got false: %s", msg)
	}
	return nil
}
