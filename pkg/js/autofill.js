function autofillPayload() {
  const checkboxes = document.querySelectorAll('input[type="checkbox"]');
  const autofillCheckboxes = Array.from(checkboxes)
    .filter((checkbox) => checkbox.id.endsWith("-autofill"))
    .filter((checkbox) => checkbox.checked);

  const checkboxPayload = autofillCheckboxes.map((item) => ({
    id: item.getAttribute("target"),
    checksum: item.getAttribute("checksum"),
    label: item.id.replace("-value-autofill", ""),
    value: asNumberIfNumeric(
      document.getElementById(item.getAttribute("target"))?.value,
    ),
  }));

  const inputFields = document.querySelectorAll("input");
  const valueFields = Array.from(inputFields).filter((f) =>
    f.id.endsWith("-value"),
  );

  const kindField = document.getElementById("type-select");
  const formPayload = { kind: kindField.value };

  for (const f of valueFields) {
    formPayload[f.id.replace("-value", "")] = asNumberIfNumeric(f.value);
  }

  return { state: formPayload, fields: checkboxPayload };
}

function applyAutofill(result) {
  if (!result || !result.data || !Array.isArray(result.data)) {
    console.error("Invalid result data for autofill");
    return;
  }

  for (const item of result.data) {
    if (!item || !item.id) {
      console.warn("Skipping invalid item in autofill data");
      continue;
    }

    const field = document.getElementById(item.id);
    if (!field) {
      console.warn(`Element with ID '${item.Id}' not found`);
      continue;
    }

    if (item.checksum !== undefined && item.checksum !== null) {
      document
        .getElementById(item.id + "-autofill")
        ?.setAttribute("checksum", item.checksum);
    }

    if (item.value !== undefined && item.value !== null) {
      field.value = item.value;
    }
  }
}

function asNumberIfNumeric(v) {
  if (v === "") return "";
  const num = Number(v);
  return isNaN(num) ? v : num;
}

async function doAutofill() {
  try {
    const payload = autofillPayload();
    const resp = await fetch("/autofill", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!resp.ok) {
      throw new Error(`HTTP error! status: ${resp.status}`);
    }

    const result = await resp.json();
    applyAutofill(result);
  } catch (error) {
    console.error("Autofill failed:", error);
    alert("Autofill failed: " + error.message);
  }
}
