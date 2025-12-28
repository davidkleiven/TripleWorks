ALTER TABLE equivalent_injections ADD CONSTRAINT fk_reactive_capability_curve FOREIGN KEY (reactive_capability_curve_id) REFERENCES reactive_capability_curves(id); 
ALTER TABLE equivalent_injections ADD CONSTRAINT fk_equivalent_equipment FOREIGN KEY (equivalent_equipment_id) REFERENCES equivalent_equipments(id); 
ALTER TABLE tie_flows ADD CONSTRAINT fk_control_area FOREIGN KEY (control_area_id) REFERENCES control_areas(id);

ALTER TABLE tie_flows ADD CONSTRAINT fk_terminal FOREIGN KEY (terminal_id) REFERENCES terminals(id);

ALTER TABLE transformer_ends ADD CONSTRAINT fk_base_voltage FOREIGN KEY (base_voltage_id) REFERENCES base_voltages(id);

ALTER TABLE transformer_ends ADD CONSTRAINT fk_terminal FOREIGN KEY (terminal_id) REFERENCES terminals(id);

ALTER TABLE transformer_ends ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE dc_terminals ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE dc_terminals ADD CONSTRAINT fk_dc_base_terminal FOREIGN KEY (dc_base_terminal_id) REFERENCES dc_base_terminals(id);

ALTER TABLE non_conform_load_schedules ADD CONSTRAINT fk_non_conform_load_group FOREIGN KEY (non_conform_load_group_id) REFERENCES non_conform_load_groups(id);

ALTER TABLE non_conform_load_schedules ADD CONSTRAINT fk_season_day_type_schedule FOREIGN KEY (season_day_type_schedule_id) REFERENCES season_day_type_schedules(id);

ALTER TABLE phase_tap_changers ADD CONSTRAINT fk_transformer_end FOREIGN KEY (transformer_end_id) REFERENCES transformer_ends(id);

ALTER TABLE phase_tap_changers ADD CONSTRAINT fk_tap_changer FOREIGN KEY (tap_changer_id) REFERENCES tap_changers(id);

ALTER TABLE generating_units ADD CONSTRAINT fk_equipment FOREIGN KEY (equipment_id) REFERENCES equipments(id);

ALTER TABLE vs_converters ADD CONSTRAINT fk_capability_curve FOREIGN KEY (capability_curve_id) REFERENCES capability_curves(id);

ALTER TABLE vs_converters ADD CONSTRAINT fk_ac_dc_converter FOREIGN KEY (ac_dc_converter_id) REFERENCES ac_dc_converters(id);

ALTER TABLE regulating_controls ADD CONSTRAINT fk_terminal FOREIGN KEY (terminal_id) REFERENCES terminals(id);

ALTER TABLE regulating_controls ADD CONSTRAINT fk_power_system_resource FOREIGN KEY (power_system_resource_id) REFERENCES power_system_resources(id);

ALTER TABLE hydro_generating_units ADD CONSTRAINT fk_hydro_power_plant FOREIGN KEY (hydro_power_plant_id) REFERENCES hydro_power_plants(id);

ALTER TABLE hydro_generating_units ADD CONSTRAINT fk_generating_unit FOREIGN KEY (generating_unit_id) REFERENCES generating_units(id);

ALTER TABLE curve_datas ADD CONSTRAINT fk_curve FOREIGN KEY (curve_id) REFERENCES curves(id);

ALTER TABLE current_limits ADD CONSTRAINT fk_operational_limit FOREIGN KEY (operational_limit_id) REFERENCES operational_limits(id);

ALTER TABLE dc_grounds ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE shunt_compensators ADD CONSTRAINT fk_regulating_cond_eq FOREIGN KEY (regulating_cond_eq_id) REFERENCES regulating_cond_eqs(id);

ALTER TABLE ac_dc_converters ADD CONSTRAINT fk_pcc_terminal FOREIGN KEY (pcc_terminal_id) REFERENCES pcc_terminals(id);

ALTER TABLE ac_dc_converters ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE wind_generating_units ADD CONSTRAINT fk_generating_unit FOREIGN KEY (generating_unit_id) REFERENCES generating_units(id);

ALTER TABLE bus_name_markers ADD CONSTRAINT fk_reporting_group FOREIGN KEY (reporting_group_id) REFERENCES reporting_groups(id);

ALTER TABLE bus_name_markers ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE equivalent_shunts ADD CONSTRAINT fk_equivalent_equipment FOREIGN KEY (equivalent_equipment_id) REFERENCES equivalent_equipments(id);

ALTER TABLE dc_line_segments ADD CONSTRAINT fk_per_length_parameter FOREIGN KEY (per_length_parameter_id) REFERENCES per_length_parameters(id);

ALTER TABLE dc_line_segments ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE load_response_characteristics ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE synchronous_machines ADD CONSTRAINT fk_initial_reactive_capability_curve FOREIGN KEY (initial_reactive_capability_curve_id) REFERENCES initial_reactive_capability_curves(id);

ALTER TABLE synchronous_machines ADD CONSTRAINT fk_rotating_machine FOREIGN KEY (rotating_machine_id) REFERENCES rotating_machines(id);

ALTER TABLE energy_sources ADD CONSTRAINT fk_energy_scheduling_type FOREIGN KEY (energy_scheduling_type_id) REFERENCES energy_scheduling_types(id);

ALTER TABLE energy_sources ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE base_voltages ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE switches ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE voltage_levels ADD CONSTRAINT fk_substation FOREIGN KEY (substation_id) REFERENCES substations(id);

ALTER TABLE voltage_levels ADD CONSTRAINT fk_base_voltage FOREIGN KEY (base_voltage_id) REFERENCES base_voltages(id);

ALTER TABLE voltage_levels ADD CONSTRAINT fk_equipment_container FOREIGN KEY (equipment_container_id) REFERENCES equipment_containers(id);

ALTER TABLE external_network_injections ADD CONSTRAINT fk_regulating_cond_eq FOREIGN KEY (regulating_cond_eq_id) REFERENCES regulating_cond_eqs(id);

ALTER TABLE phase_tap_changer_asymmetricals ADD CONSTRAINT fk_phase_tap_changer_non_linear FOREIGN KEY (phase_tap_changer_non_linear_id) REFERENCES phase_tap_changer_non_linears(id);

ALTER TABLE control_area_generating_units ADD CONSTRAINT fk_control_area FOREIGN KEY (control_area_id) REFERENCES control_areas(id);

ALTER TABLE control_area_generating_units ADD CONSTRAINT fk_generating_unit FOREIGN KEY (generating_unit_id) REFERENCES generating_units(id);

ALTER TABLE control_area_generating_units ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE conform_loads ADD CONSTRAINT fk_load_group FOREIGN KEY (load_group_id) REFERENCES load_groups(id);

ALTER TABLE conform_loads ADD CONSTRAINT fk_energy_consumer FOREIGN KEY (energy_consumer_id) REFERENCES energy_consumers(id);

ALTER TABLE basic_interval_schedules ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE nonlinear_shunt_compensator_points ADD CONSTRAINT fk_nonlinear_shunt_compensator FOREIGN KEY (nonlinear_shunt_compensator_id) REFERENCES nonlinear_shunt_compensators(id);

ALTER TABLE phase_tap_changer_non_linears ADD CONSTRAINT fk_phase_tap_changer FOREIGN KEY (phase_tap_changer_id) REFERENCES phase_tap_changers(id);

ALTER TABLE phase_tap_changer_tabulars ADD CONSTRAINT fk_phase_tap_changer_table FOREIGN KEY (phase_tap_changer_table_id) REFERENCES phase_tap_changer_tables(id);

ALTER TABLE phase_tap_changer_tabulars ADD CONSTRAINT fk_phase_tap_changer FOREIGN KEY (phase_tap_changer_id) REFERENCES phase_tap_changers(id);

ALTER TABLE voltage_limits ADD CONSTRAINT fk_operational_limit FOREIGN KEY (operational_limit_id) REFERENCES operational_limits(id);

ALTER TABLE operational_limit_types ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE linear_shunt_compensators ADD CONSTRAINT fk_shunt_compensator FOREIGN KEY (shunt_compensator_id) REFERENCES shunt_compensators(id);

ALTER TABLE ac_dc_converter_dc_terminals ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE ac_dc_converter_dc_terminals ADD CONSTRAINT fk_dc_base_terminal FOREIGN KEY (dc_base_terminal_id) REFERENCES dc_base_terminals(id);

ALTER TABLE ac_dc_terminals ADD CONSTRAINT fk_bus_name_marker FOREIGN KEY (bus_name_marker_id) REFERENCES bus_name_markers(id);

ALTER TABLE ac_dc_terminals ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE sub_geographical_regions ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id);

ALTER TABLE sub_geographical_regions ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE hydro_power_plants ADD CONSTRAINT fk_power_system_resource FOREIGN KEY (power_system_resource_id) REFERENCES power_system_resources(id);

ALTER TABLE operational_limits ADD CONSTRAINT fk_operational_limit_set FOREIGN KEY (operational_limit_set_id) REFERENCES operational_limit_sets(id);

ALTER TABLE operational_limits ADD CONSTRAINT fk_operational_limit_type FOREIGN KEY (operational_limit_type_id) REFERENCES operational_limit_types(id);

ALTER TABLE operational_limits ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE equivalent_branchs ADD CONSTRAINT fk_equivalent_equipment FOREIGN KEY (equivalent_equipment_id) REFERENCES equivalent_equipments(id);

ALTER TABLE ratio_tap_changer_table_points ADD CONSTRAINT fk_ratio_tap_changer_table FOREIGN KEY (ratio_tap_changer_table_id) REFERENCES ratio_tap_changer_tables(id);

ALTER TABLE ratio_tap_changer_table_points ADD CONSTRAINT fk_tap_changer_table_point FOREIGN KEY (tap_changer_table_point_id) REFERENCES tap_changer_table_points(id);

ALTER TABLE dc_nodes ADD CONSTRAINT fk_dc_equipment_container FOREIGN KEY (dc_equipment_container_id) REFERENCES dc_equipment_containers(id);

ALTER TABLE dc_nodes ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE terminals ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE terminals ADD CONSTRAINT fk_ac_dc_terminal FOREIGN KEY (ac_dc_terminal_id) REFERENCES ac_dc_terminals(id);

ALTER TABLE dc_base_terminals ADD CONSTRAINT fk_dc_node FOREIGN KEY (dc_node_id) REFERENCES dc_nodes(id);

ALTER TABLE dc_base_terminals ADD CONSTRAINT fk_ac_dc_terminal FOREIGN KEY (ac_dc_terminal_id) REFERENCES ac_dc_terminals(id);

ALTER TABLE energy_consumers ADD CONSTRAINT fk_load_response FOREIGN KEY (load_response_id) REFERENCES load_responses(id);

ALTER TABLE energy_consumers ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE conform_load_schedules ADD CONSTRAINT fk_conform_load_group FOREIGN KEY (conform_load_group_id) REFERENCES conform_load_groups(id);

ALTER TABLE conform_load_schedules ADD CONSTRAINT fk_season_day_type_schedule FOREIGN KEY (season_day_type_schedule_id) REFERENCES season_day_type_schedules(id);

ALTER TABLE lines ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id);

ALTER TABLE lines ADD CONSTRAINT fk_equipment_container FOREIGN KEY (equipment_container_id) REFERENCES equipment_containers(id);

ALTER TABLE equivalent_equipments ADD CONSTRAINT fk_equivalent_network FOREIGN KEY (equivalent_network_id) REFERENCES equivalent_networks(id);

ALTER TABLE equivalent_equipments ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE dc_lines ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id);

ALTER TABLE dc_lines ADD CONSTRAINT fk_dc_equipment_container FOREIGN KEY (dc_equipment_container_id) REFERENCES dc_equipment_containers(id);

ALTER TABLE phase_tap_changer_table_points ADD CONSTRAINT fk_phase_tap_changer_table FOREIGN KEY (phase_tap_changer_table_id) REFERENCES phase_tap_changer_tables(id);

ALTER TABLE phase_tap_changer_table_points ADD CONSTRAINT fk_tap_changer_table_point FOREIGN KEY (tap_changer_table_point_id) REFERENCES tap_changer_table_points(id);

ALTER TABLE regulating_cond_eqs ADD CONSTRAINT fk_regulating_control FOREIGN KEY (regulating_control_id) REFERENCES regulating_controls(id);

ALTER TABLE regulating_cond_eqs ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE phase_tap_changer_linears ADD CONSTRAINT fk_phase_tap_changer FOREIGN KEY (phase_tap_changer_id) REFERENCES phase_tap_changers(id);

ALTER TABLE cs_converters ADD CONSTRAINT fk_ac_dc_converter FOREIGN KEY (ac_dc_converter_id) REFERENCES ac_dc_converters(id);

ALTER TABLE operational_limit_sets ADD CONSTRAINT fk_terminal FOREIGN KEY (terminal_id) REFERENCES terminals(id);

ALTER TABLE operational_limit_sets ADD CONSTRAINT fk_equipment FOREIGN KEY (equipment_id) REFERENCES equipments(id);

ALTER TABLE operational_limit_sets ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE ratio_tap_changers ADD CONSTRAINT fk_ratio_tap_changer_table FOREIGN KEY (ratio_tap_changer_table_id) REFERENCES ratio_tap_changer_tables(id);

ALTER TABLE ratio_tap_changers ADD CONSTRAINT fk_transformer_end FOREIGN KEY (transformer_end_id) REFERENCES transformer_ends(id);

ALTER TABLE ratio_tap_changers ADD CONSTRAINT fk_tap_changer FOREIGN KEY (tap_changer_id) REFERENCES tap_changers(id);

ALTER TABLE conducting_equipments ADD CONSTRAINT fk_base_voltage FOREIGN KEY (base_voltage_id) REFERENCES base_voltages(id);

ALTER TABLE conducting_equipments ADD CONSTRAINT fk_equipment FOREIGN KEY (equipment_id) REFERENCES equipments(id);

ALTER TABLE series_compensators ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE power_transformer_ends ADD CONSTRAINT fk_power_transformer FOREIGN KEY (power_transformer_id) REFERENCES power_transformers(id);

ALTER TABLE power_transformer_ends ADD CONSTRAINT fk_transformer_end FOREIGN KEY (transformer_end_id) REFERENCES transformer_ends(id);

ALTER TABLE static_var_compensators ADD CONSTRAINT fk_regulating_cond_eq FOREIGN KEY (regulating_cond_eq_id) REFERENCES regulating_cond_eqs(id);

ALTER TABLE non_conform_loads ADD CONSTRAINT fk_load_group FOREIGN KEY (load_group_id) REFERENCES load_groups(id);

ALTER TABLE non_conform_loads ADD CONSTRAINT fk_energy_consumer FOREIGN KEY (energy_consumer_id) REFERENCES energy_consumers(id);

ALTER TABLE equipments ADD CONSTRAINT fk_equipment_container FOREIGN KEY (equipment_container_id) REFERENCES equipment_containers(id);

ALTER TABLE equipments ADD CONSTRAINT fk_power_system_resource FOREIGN KEY (power_system_resource_id) REFERENCES power_system_resources(id);

ALTER TABLE control_areas ADD CONSTRAINT fk_power_system_resource FOREIGN KEY (power_system_resource_id) REFERENCES power_system_resources(id);

ALTER TABLE curves ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE asynchronous_machines ADD CONSTRAINT fk_rotating_machine FOREIGN KEY (rotating_machine_id) REFERENCES rotating_machines(id);

ALTER TABLE dc_shunts ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE substations ADD CONSTRAINT fk_region FOREIGN KEY (region_id) REFERENCES regions(id);

ALTER TABLE substations ADD CONSTRAINT fk_equipment_container FOREIGN KEY (equipment_container_id) REFERENCES equipment_containers(id);

ALTER TABLE conductors ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE fossil_fuels ADD CONSTRAINT fk_thermal_generating_unit FOREIGN KEY (thermal_generating_unit_id) REFERENCES thermal_generating_units(id);

ALTER TABLE fossil_fuels ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE dc_seriesdevices ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE tap_changers ADD CONSTRAINT fk_tap_changer_control FOREIGN KEY (tap_changer_control_id) REFERENCES tap_changer_controls(id);

ALTER TABLE tap_changers ADD CONSTRAINT fk_power_system_resource FOREIGN KEY (power_system_resource_id) REFERENCES power_system_resources(id);

ALTER TABLE regular_interval_schedules ADD CONSTRAINT fk_basic_interval_schedule FOREIGN KEY (basic_interval_schedule_id) REFERENCES basic_interval_schedules(id);

ALTER TABLE rotating_machines ADD CONSTRAINT fk_generating_unit FOREIGN KEY (generating_unit_id) REFERENCES generating_units(id);

ALTER TABLE rotating_machines ADD CONSTRAINT fk_regulating_cond_eq FOREIGN KEY (regulating_cond_eq_id) REFERENCES regulating_cond_eqs(id);

ALTER TABLE dc_converterunits ADD CONSTRAINT fk_substation FOREIGN KEY (substation_id) REFERENCES substations(id);

ALTER TABLE dc_converterunits ADD CONSTRAINT fk_dc_equipment_container FOREIGN KEY (dc_equipment_container_id) REFERENCES dc_equipment_containers(id);

ALTER TABLE hydro_pumps ADD CONSTRAINT fk_rotating_machine FOREIGN KEY (rotating_machine_id) REFERENCES rotating_machines(id);

ALTER TABLE hydro_pumps ADD CONSTRAINT fk_hydro_power_plant FOREIGN KEY (hydro_power_plant_id) REFERENCES hydro_power_plants(id);

ALTER TABLE hydro_pumps ADD CONSTRAINT fk_equipment FOREIGN KEY (equipment_id) REFERENCES equipments(id);

ALTER TABLE ac_linesegments ADD CONSTRAINT fk_conductor FOREIGN KEY (conductor_id) REFERENCES conductors(id);

ALTER TABLE busbar_sections ADD CONSTRAINT fk_connector FOREIGN KEY (connector_id) REFERENCES connectors(id);

ALTER TABLE connectors ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

ALTER TABLE equipment_containers ADD CONSTRAINT fk_connectivity_node_container FOREIGN KEY (connectivity_node_container_id) REFERENCES connectivity_node_containers(id);

ALTER TABLE power_system_resources ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE nonlinear_shunt_compensators ADD CONSTRAINT fk_shunt_compensator FOREIGN KEY (shunt_compensator_id) REFERENCES shunt_compensators(id);

ALTER TABLE load_groups ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE dc_conducting_equipments ADD CONSTRAINT fk_equipment FOREIGN KEY (equipment_id) REFERENCES equipments(id);

ALTER TABLE equivalent_networks ADD CONSTRAINT fk_connectivity_node_container FOREIGN KEY (connectivity_node_container_id) REFERENCES connectivity_node_containers(id);

ALTER TABLE connectivity_node_containers ADD CONSTRAINT fk_power_system_resource FOREIGN KEY (power_system_resource_id) REFERENCES power_system_resources(id);

ALTER TABLE dc_equipment_containers ADD CONSTRAINT fk_equipment_container FOREIGN KEY (equipment_container_id) REFERENCES equipment_containers(id);

ALTER TABLE solar_generating_units ADD CONSTRAINT fk_generating_unit FOREIGN KEY (generating_unit_id) REFERENCES generating_units(id);

ALTER TABLE ratio_tap_changer_tables ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE non_conform_load_groups ADD CONSTRAINT fk_load_group FOREIGN KEY (load_group_id) REFERENCES load_groups(id);

ALTER TABLE energy_scheduling_types ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE protected_switches ADD CONSTRAINT fk_switch FOREIGN KEY (switch_id) REFERENCES switches(id);

ALTER TABLE disconnectors ADD CONSTRAINT fk_switch FOREIGN KEY (switch_id) REFERENCES switches(id);

ALTER TABLE phase_tap_changer_symmetricals ADD CONSTRAINT fk_phase_tap_changer_non_linear FOREIGN KEY (phase_tap_changer_non_linear_id) REFERENCES phase_tap_changer_non_linears(id);

ALTER TABLE dc_switches ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE junctions ADD CONSTRAINT fk_connector FOREIGN KEY (connector_id) REFERENCES connectors(id);

ALTER TABLE reactive_capability_curves ADD CONSTRAINT fk_curve FOREIGN KEY (curve_id) REFERENCES curves(id);

ALTER TABLE phase_tap_changer_tables ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE tap_changer_controls ADD CONSTRAINT fk_regulating_control FOREIGN KEY (regulating_control_id) REFERENCES regulating_controls(id);

ALTER TABLE thermal_generating_units ADD CONSTRAINT fk_generating_unit FOREIGN KEY (generating_unit_id) REFERENCES generating_units(id);

ALTER TABLE dc_disconnectors ADD CONSTRAINT fk_dc_switch FOREIGN KEY (dc_switch_id) REFERENCES dc_switches(id);

ALTER TABLE reporting_groups ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE dc_choppers ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE breakers ADD CONSTRAINT fk_protected_switch FOREIGN KEY (protected_switch_id) REFERENCES protected_switches(id);

ALTER TABLE geographical_regions ADD CONSTRAINT fk_identified_object FOREIGN KEY (identified_object_id) REFERENCES identified_objects(id);

ALTER TABLE load_break_switches ADD CONSTRAINT fk_protected_switch FOREIGN KEY (protected_switch_id) REFERENCES protected_switches(id);

ALTER TABLE conform_load_groups ADD CONSTRAINT fk_load_group FOREIGN KEY (load_group_id) REFERENCES load_groups(id);

ALTER TABLE dc_breakers ADD CONSTRAINT fk_dc_switch FOREIGN KEY (dc_switch_id) REFERENCES dc_switches(id);

ALTER TABLE dc_busbars ADD CONSTRAINT fk_dc_conducting_equipment FOREIGN KEY (dc_conducting_equipment_id) REFERENCES dc_conducting_equipments(id);

ALTER TABLE nuclear_generating_units ADD CONSTRAINT fk_generating_unit FOREIGN KEY (generating_unit_id) REFERENCES generating_units(id);

ALTER TABLE vs_capability_curves ADD CONSTRAINT fk_curve FOREIGN KEY (curve_id) REFERENCES curves(id);

ALTER TABLE power_transformers ADD CONSTRAINT fk_conducting_equipment FOREIGN KEY (conducting_equipment_id) REFERENCES conducting_equipments(id);

