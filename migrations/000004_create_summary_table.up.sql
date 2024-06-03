CREATE TABLE IF NOT EXISTS summary (
   id SERIAL ,
   owner_id UUID NOT NULL ,
   skills TEXT NOT NULL ,
   bio TEXT NOT NULL ,
   languages TEXT NOT NULL
);

INSERT INTO summary (owner_id, skills, bio, languages) VALUES
('9362646a-594e-4ff8-b9e5-96060faba9f5', 'Skill 1', 'Bio 1', 'Language 1'),
('25846404-2939-40bc-98ee-4654252b411e', 'Skill 2', 'Bio 2', 'Language 2'),
('e08652ed-9146-4d3f-b19a-5b7f6e3a48fa', 'Skill 3', 'Bio 3', 'Language 3'),
('9362646a-594e-4ff8-b9e5-96060faba9f5', 'Skill 4', 'Bio 4', 'Language 4'),
('25846404-2939-40bc-98ee-4654252b411e', 'Skill 5', 'Bio 5', 'Language 5');


