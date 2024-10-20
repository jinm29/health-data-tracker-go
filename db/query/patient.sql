-- name: CreatePatient :one
INSERT INTO patients (full_name, age, gender, contact_info, assigned_device_id)
Values ($1, $2, $3, $4, $5)
RETURNING id, full_name, assigned_device_id;

-- name: GetAllPatients :many
SELECT * FROM patients;

-- name: GetPatientWithDevice :one
SELECT 
    p.id AS patient_id,
    p.full_name,
    p.age,
    p.gender,
    p.contact_info,
    wd.id AS device_id,
    wd.device_name,
    wd.device_type
FROM 
    patients p
INNER JOIN 
    wearable_devices wd ON p.assigned_device_id = wd.id
WHERE 
    p.id = $1;

-- name: AssignDeviceToPatient :exec
UPDATE patients
SET assigned_device_id = $1
WHERE id = $2;
