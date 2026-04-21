select
    t.mrid as terminal_mrid,
    t.name as terminal_name,
    t.sequence_number as terminal_sequence_number,
    c.mrid as connectivity_node_mrid,
    c.name as connectivity_node_name,
    v.mrid as voltage_level_mrid,
    v.name as voltage_level_name,
    s.mrid as substation_mrid,
    s.name as substation_name
from v_terminals_latest t
inner join v_connectivity_nodes_latest c on t.connectivity_node_mrid = c.mrid
inner join
    v_voltage_levels_latest v
    on c.connectivity_node_container_mrid = v.mrid
inner join v_substations_latest s on v.substation_mrid = s.mrid
where t.conducting_equipment_mrid = ?
order by s.mrid
