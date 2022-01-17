--* get what a teacher teaches
SELECT insegna.ID, c.materia, c.api_value, i.nome, i.cognome, i.matricola FROM insegna join corsi c on id_corso=c.ID join insegnante i on id_insegnante=i.ID;
-- ID,materia,api_value,nome,cognome,matricola

--* get the subjet and the api value of a course
SELECT materia, api_value FROM corsi;
-- materia,api_value

--* get the exam informations
SELECT difficolta, numero_domande, nome, c.materia, c.api_value FROM esami e JOIN corsi c on ID_corso = c.ID;
-- difficolta,numero_domande,nome,materia,api_value

--* get the result of an exam made by a student
SELECT r.nome_studente, r.cognome_studente, r.contenuto, r.tentativi, e.difficolta, e.nome as nome_esame, e.numero_domande, c.materia, i.matricola as matricola_insegnante FROM risultati_esami r join esami e on ID_esame=e.ID join corsi c on e.ID_corso = c.ID join insegnante i on ID_insegnante = i.ID;
-- nome_studente,cognome_studente,contenuto,tentativi,difficolta,nome_esame,numero_domande,materia,matricola_insegnante

--* get the information of a exam to make the api request to trivia db
SELECT e.difficolta, e.numero_domande, e.nome, c.api_value FROM esami e join corsi c on ID_corso = c.ID;
-- difficolta,numero_domande,nome,api_value




