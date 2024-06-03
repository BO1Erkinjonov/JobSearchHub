CREATE TABLE IF NOT EXISTS requests (
    job_id UUID NOT NULL,
    client_id UUID NOT NULL,
    summary_id INT NOT NULL,
    status_resp VARCHAR(100) CHECK (status_resp IN ('refusal', 'in expectation', 'accepted')) DEFAULT 'in expectation' NOT NULL,
    description_resp TEXT
);

INSERT INTO requests (job_id, client_id, summary_id, status_resp, description_resp) VALUES
('5e27c090-bdd4-4c84-b916-70eb160d89bc', '9362646a-594e-4ff8-b9e5-96060faba9f5', 1, 'in expectation', 'Description 1'),
('c1e2d969-a837-4cd8-b9ab-5bf7dddb1558', '25846404-2939-40bc-98ee-4654252b411e', 2, 'in expectation', 'Description 2'),
('1f79c7d6-81ab-4044-a503-25884b1cbb71', 'e08652ed-9146-4d3f-b19a-5b7f6e3a48fa', 3, 'in expectation', 'Description 3'),
('c1e2d969-a837-4cd8-b9ab-5bf7dddb1558', '9362646a-594e-4ff8-b9e5-96060faba9f5', 4, 'accepted', 'Description 4'),
('071dbb24-2b48-464e-b671-1b8716b238a1', '25846404-2939-40bc-98ee-4654252b411e', 5, 'accepted', 'Description 5'),
('0814fc16-8762-4784-a665-72422833296c', 'e08652ed-9146-4d3f-b19a-5b7f6e3a48fa', 3, 'accepted', 'Description 6'),
('5e27c090-bdd4-4c84-b916-70eb160d89bc', '9362646a-594e-4ff8-b9e5-96060faba9f5', 4, 'refusal', 'Description 7'),
('5e27c090-bdd4-4c84-b916-70eb160d89bc', '25846404-2939-40bc-98ee-4654252b411e', 5, 'refusal', 'Description 8'),
('0814fc16-8762-4784-a665-72422833296c', 'e08652ed-9146-4d3f-b19a-5b7f6e3a48fa', 3, 'refusal', 'Description 9'),
('f4da6c9a-fbc4-4a12-9fd1-9944cfde6692', '9362646a-594e-4ff8-b9e5-96060faba9f5', 1, 'in expectation', 'Description 10');


