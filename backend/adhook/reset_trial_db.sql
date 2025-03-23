DO $$DECLARE
    r RECORD;
BEGIN
    -- 全テーブルのデータを削除 (外部キー制約がある場合は順序に注意)
FOR r IN (
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
        AND tablename IN (
            'departments', 'facilities', 'facility_shared_links', 'gantt_groups', 'holidays',
            'milestones', 'operation_settings', 'processes', 'simulation_departments',
            'simulation_facilities', 'simulation_facility_shared_links', 'simulation_gantt_groups',
            'simulation_holidays', 'simulation_locks', 'simulation_milestones',
            'simulation_operation_settings', 'simulation_processes', 'simulation_ticket_users',
            'simulation_tickets', 'simulation_units', 'simulation_users', 'ticket_users',
            'tickets', 'units', 'users'
        )
    ) LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE;';
END LOOP;

    -- 全シーケンスをリセット (TRUNCATE ... RESTART IDENTITYで自動的にリセットされるけど、念のため)
FOR r IN (
        SELECT sequencename
        FROM pg_sequences
        WHERE schemaname = 'public'
        AND sequencename IN (
            'departments_id_seq', 'facilities_id_seq', 'facility_shared_links_id_seq',
            'gantt_groups_id_seq', 'holidays_id_seq', 'milestones_id_seq',
            'operation_settings_id_seq', 'processes_id_seq', 'simulation_departments_id_seq',
            'simulation_facilities_id_seq', 'simulation_facility_shared_links_id_seq',
            'simulation_gantt_groups_id_seq', 'simulation_holidays_id_seq',
            'simulation_milestones_id_seq', 'simulation_operation_settings_id_seq',
            'simulation_processes_id_seq', 'simulation_ticket_users_id_seq',
            'simulation_tickets_id_seq', 'simulation_units_id_seq', 'simulation_users_id_seq',
            'ticket_users_id_seq', 'tickets_id_seq', 'units_id_seq', 'users_id_seq'
        )
    ) LOOP
        EXECUTE 'ALTER SEQUENCE ' || quote_ident(r.sequencename) || ' RESTART WITH 1;';
END LOOP;
END$$;

INSERT INTO public.users (id, department_id, limit_of_operation, last_name, first_name, password, email, role, created_at, updated_at)
VALUES (1, 0, 0, '管理者', 'システム', '$2a$10$/zqYDZLTELr/bXRg7ROEJOt8E8hxNmfyPVeg9VwXVocDgHfl6OrLW', 'admin', 'admin', '2024-12-05 03:43:37.033748 +00:00', 1733370217);
-- admin/itumono