create view v_substation_bidzones_latest as
select s.mrid, s.name, s.region_mrid, g.name as bidzone
from v_substations_latest s
inner join v_sub_geographical_regions_latest g on s.region_mrid = g.mrid
