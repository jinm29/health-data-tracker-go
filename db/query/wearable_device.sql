-- name: AddWearableDevices :one
INSERT INTO wearable_devices (device_type, device_name, assigned_patient_id, serial_number, country_code)
Values ($1, $2, $3, $4, $5)
RETURNING id, device_type, assigned_patient_id;

-- name: GetAvailableDevices :many
SELECT * FROM wearable_devices wd WHERE wd.assigned_patient_id IS NULL;

-- name: AssignDevice :exec
UPDATE wearable_devices
SET assigned_patient_id = $1
WHERE id = $2;
