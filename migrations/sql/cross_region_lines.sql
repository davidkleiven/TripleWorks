CREATE VIEW v_cross_region_lines_latest AS
SELECT
    line_mrid, line_name,
    MAX(CASE WHEN sequence_number = 1 THEN bidzone END) AS from_bidzone,
    MAX(CASE WHEN sequence_number = 2 THEN bidzone END) AS to_bidzone
FROM (
    SELECT
        a.mrid AS line_mrid,
        a.name AS line_name,
        sg.name AS bidzone,
        t.sequence_number
    FROM v_ac_line_segments_latest a
    INNER JOIN v_terminals_latest t ON t.conducting_equipment_mrid = a.mrid
    INNER JOIN
        v_connectivity_nodes_latest c
        ON t.connectivity_node_mrid = c.mrid
    INNER JOIN
        v_voltage_levels_latest v
        ON c.connectivity_node_container_mrid = v.mrid
    INNER JOIN v_substations_latest s ON v.substation_mrid = s.mrid
    INNER JOIN v_sub_geographical_regions_latest sg ON sg.mrid = s.region_mrid
) AS lines
GROUP BY line_mrid, line_name
HAVING COUNT(DISTINCT bidzone) = 2;
