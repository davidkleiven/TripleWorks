
create view v_invalid_lines as
with a_mrid_t_mrid as (
    select a.mrid as a_mrid, t.mrid as t_mrid, t.sequence_number as seq_no
    from v_ac_line_segments_latest a
    inner join v_terminals_latest t on t.conducting_equipment_mrid = a.mrid
),

invalid_lines as (
    select a_mrid
    from a_mrid_t_mrid aa
    group by a_mrid
    having (count(*) != 2 or count(distinct aa.seq_no) != 2)
)

select
    a.mrid as line_mrid,
    a.name as line_name,
    s.mrid as sub_mrid,
    s.name as sub_name,
    t.mrid as terminal_mrid,
    t.sequence_number as terminal_sequence_number
from v_ac_line_segments_latest a
inner join v_terminals_latest t on t.conducting_equipment_mrid = a.mrid
inner join v_connectivity_nodes_latest c on t.connectivity_node_mrid = c.mrid
inner join
    v_voltage_levels_latest v
    on c.connectivity_node_container_mrid = v.mrid
inner join v_substations_latest s on v.substation_mrid = s.mrid
where a.mrid in (select a_mrid from invalid_lines)
order by line_mrid
