CREATE TABLE IF NOT EXISTS jobs (
    id UUID NOT NULL,
    owner_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    responses INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);


INSERT INTO jobs (id, owner_id, title, description, responses, created_at, updated_at, deleted_at) VALUES
('5e27c090-bdd4-4c84-b916-70eb160d89bc', 'e08652ed-9146-4d3f-b19a-5b7f6e3a48fa', 'Fake Job 1', 'Description 1', 3, CURRENT_TIMESTAMP, NULL, NULL),
('1f79c7d6-81ab-4044-a503-25884b1cbb71', '25846404-2939-40bc-98ee-4654252b411e', 'Fake Job 2', 'Description 2', 1, CURRENT_TIMESTAMP, NULL, NULL),
('11473b10-9fc7-4415-af03-b65aee6b0c91', 'fa46bb60-546b-4cee-8404-9924d14f0a2d', 'Fake Job 3', 'Description 3', 0, CURRENT_TIMESTAMP, NULL, NULL),
('c1e2d969-a837-4cd8-b9ab-5bf7dddb1558', 'd92b5489-1055-45cc-89fa-6f7fb4a44000', 'Fake Job 4', 'Description 4', 2, CURRENT_TIMESTAMP, NULL, NULL),
('f4da6c9a-fbc4-4a12-9fd1-9944cfde6692', '1f993884-8b16-4a6d-87bc-80aaf6ddb97c', 'Fake Job 5', 'Description 5', 1, CURRENT_TIMESTAMP, NULL, NULL),
('071dbb24-2b48-464e-b671-1b8716b238a1', 'e08652ed-9146-4d3f-b19a-5b7f6e3a48fa', 'Fake Job 6', 'Description 6', 1, CURRENT_TIMESTAMP, NULL, NULL),
('0814fc16-8762-4784-a665-72422833296c', '25846404-2939-40bc-98ee-4654252b411e', 'Fake Job 7', 'Description 7', 2, CURRENT_TIMESTAMP, NULL, NULL),
('e457807f-4537-490c-b840-2815ad55b543', 'fa46bb60-546b-4cee-8404-9924d14f0a2d', 'Fake Job 8', 'Description 8', 0, CURRENT_TIMESTAMP, NULL, NULL),
('1b189529-3563-41ba-af24-23dc41622b8f', 'd92b5489-1055-45cc-89fa-6f7fb4a44000', 'Fake Job 9', 'Description 9', 0, CURRENT_TIMESTAMP, NULL, NULL),
('c79dd4c8-94b5-4bea-ba60-b99394ea1e09', '1f993884-8b16-4a6d-87bc-80aaf6ddb97c', 'Fake Job 10', 'Description 10', 0, CURRENT_TIMESTAMP, NULL, NULL);


