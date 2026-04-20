WITH terminal_substation AS (
    SELECT t.mrid, v.substation_mrid FROM v_terminals_latest t
    INNER JOIN
        v_connectivity_nodes_latest c
        ON t.connectivity_node_mrid = c.mrid
    INNER JOIN
        v_voltage_levels_latest v
        ON c.connectivity_node_container_mrid = v.mrid
)

SELECT
    l.mrid, l.r, l.x, b.nominal_voltage, ts.substation_mrid, t.sequence_number
FROM v_ac_line_segments_latest l
INNER JOIN v_base_voltages_latest b ON l.base_voltage_mrid = b.mrid
INNER JOIN v_terminals_latest t ON t.conducting_equipment_mrid = l.mrid
INNER JOIN terminal_substation ts ON ts.mrid = t.mrid
