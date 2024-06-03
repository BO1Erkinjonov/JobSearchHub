CREATE TABLE IF NOT EXISTS clients (
   id UUID NOT NULL,
   role VARCHAR(64) NOT NULL ,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL,
    password VARCHAR(64) NOT NULL,
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
    );


INSERT INTO clients (id, role, first_name, last_name, email, password, refresh_token, created_at, updated_at, deleted_at) VALUES
      ('e08652ed-9146-4d3f-b19a-5b7f6e3a48fa', 'client', 'John', 'Doe', 'johndoe@gmail.com', 'password123', 'sample_refresh_token', CURRENT_TIMESTAMP, NULL, NULL),
      ('25846404-2939-40bc-98ee-4654252b411e', 'client', 'Jane', 'Smith', 'janesmith@gmail.com', 'pass123', 'another_refresh_token', CURRENT_TIMESTAMP, NULL, NULL),
      ('ec4b075f-5cd4-440f-adcd-14613eee2eed', 'admin', 'Alice', 'Johnson', 'alicejohnson@gmail.com', 'qwerty', 'refresh_me', CURRENT_TIMESTAMP, NULL, NULL),
      ('fa46bb60-546b-4cee-8404-9924d14f0a2d', 'client', 'Bob', 'Brown', 'bobbrown@gmail.com', 'secret', 'token_refresh', CURRENT_TIMESTAMP, NULL, NULL),
      ('b9ce2b4a-29d2-4940-9fc5-ef415eb75abe', 'client', 'Emily', 'Davis', 'emilydavis@gmail.com', 'letmein', 'token_refreshing', CURRENT_TIMESTAMP, NULL, NULL),
      ('c337a35a-4953-4c0b-8226-690a0f0b1f84', 'admin', 'Michael', 'Wilson', 'michaelwilson@gmail.com', 'password123', 'refresh_token_sample', CURRENT_TIMESTAMP, NULL, NULL),
      ('d92b5489-1055-45cc-89fa-6f7fb4a44000', 'client', 'Emma', 'Martinez', 'emmamartinez@gmail.com', 'pass123', 'sample_refreshing_token', CURRENT_TIMESTAMP, NULL, NULL),
      ('1f993884-8b16-4a6d-87bc-80aaf6ddb97c', 'client', 'David', 'Taylor', 'davidtaylor@gmail.com', 'qwerty', 'refresh_this_token', CURRENT_TIMESTAMP, NULL, NULL),
      ('13a8edb7-dee0-4542-a22f-574089a0ee9a', 'admin', 'Olivia', 'Anderson', 'oliviaanderson@gmail.com', 'secret', 'another_refresh_sample', CURRENT_TIMESTAMP, NULL, NULL),
      ('9362646a-594e-4ff8-b9e5-96060faba9f5', 'client', 'James', 'Thomas', 'jamesthomas@gmail.com', 'letmein', 'token_refreshing_sample', CURRENT_TIMESTAMP, NULL, NULL);

