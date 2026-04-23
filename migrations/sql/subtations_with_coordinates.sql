create view substations_geo as
select s.mrid, s.name, p.xposition as longitude, p.yposition as latitude
from v_substations_latest s
inner join v_locations_latest l on l.mrid = s.location_mrid
inner join position_points p on p.location_mrid = l.mrid
