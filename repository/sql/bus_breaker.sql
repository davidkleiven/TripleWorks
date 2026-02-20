WITH latest_terminals AS (
	SELECT mrid, connectivity_node_mrid, conducting_equipment_mrid FROM (
		SELECT *, ROW_NUMBER() OVER (PARTITION BY mrid ORDER BY commit_id DESC) AS rn FROM terminals
	) WHERE rn = 1
),
latest_con_nodes AS (
	SELECT mrid, connectivity_node_container_mrid FROM (
		SELECT *, ROW_NUMBER() OVER (PARTITION BY mrid ORDER BY commit_id DESC) AS rn FROM connectivity_nodes
	) WHERE rn = 1
),
latest_vls AS (
	SELECT mrid, substation_mrid FROM (
		SELECT *, ROW_NUMBER() OVER (PARTITION BY mrid ORDER BY commit_id DESC) AS rn FROM voltage_levels
	) WHERE rn = 1
),
latest_lines AS (
	SELECT mrid, r, x, base_voltage_mrid FROM (
		SELECT *, ROW_NUMBER() OVER (PARTITION BY mrid ORDER BY commit_id DESC) AS rn FROM ac_line_segments
	) WHERE rn = 1
),
latest_bv AS (
	SELECT mrid, nominal_voltage FROM (
		SELECT *, ROW_NUMBER() OVER (PARTITION BY mrid ORDER BY commit_id DESC) AS rn FROM base_voltages
	) WHERE rn = 1
),
terminal_substation AS (
	SELECT t.mrid, v.substation_mrid FROM latest_terminals t
	INNER JOIN latest_con_nodes c ON t.connectivity_node_mrid = c.mrid
	INNER JOIN latest_vls v ON c.connectivity_node_container_mrid = v.mrid
)
SELECT l.mrid, l.r, l.x, b.nominal_voltage, ts.substation_mrid
FROM latest_lines l
INNER JOIN latest_bv b ON l.base_voltage_mrid = b.mrid
INNER JOIN latest_terminals t ON t.conducting_equipment_mrid = l.mrid
INNER JOIN terminal_substation ts ON ts.mrid = t.mrid
