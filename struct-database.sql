-- MySQL dump 10.13  Distrib 8.0.27, for Linux (x86_64)
--
-- Host: localhost    Database: vano-teachers
-- ------------------------------------------------------
-- Server version	8.0.27-0ubuntu0.21.10.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `corsi`
--

DROP TABLE IF EXISTS `corsi`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `corsi` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `materia` varchar(45) NOT NULL,
  `api_value` int NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `corsi`
--

LOCK TABLES `corsi` WRITE;
/*!40000 ALTER TABLE `corsi` DISABLE KEYS */;
INSERT INTO `corsi` VALUES (1,'informatica',18),(2,'matematica',19),(3,'storia',23),(5,'anime e manga',31);
/*!40000 ALTER TABLE `corsi` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `esami`
--

DROP TABLE IF EXISTS `esami`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `esami` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `difficolta` varchar(10) NOT NULL,
  `ID_corso` int DEFAULT NULL,
  `numero_domande` int NOT NULL,
  `nome` varchar(45) NOT NULL,
  `ID_insegnante` int NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `ID_corso_idx` (`ID_corso`),
  KEY `ID_insegnante_idx` (`ID_insegnante`),
  CONSTRAINT `ID_corso` FOREIGN KEY (`ID_corso`) REFERENCES `corsi` (`ID`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `esami`
--

LOCK TABLES `esami` WRITE;
/*!40000 ALTER TABLE `esami` DISABLE KEYS */;
INSERT INTO `esami` VALUES (1,'easy',1,10,'informatica 101',2),(2,'hard',1,10,'goku bello goku bello',2),(3,'hard',2,10,'aaaaaa',9),(4,'hard',1,12,'aaa',2),(5,'easy',5,10,'yuri and love',2),(6,'hard',5,20,'vanella e gli anime',11),(7,'hard',2,2,'Analisi 1',12);
/*!40000 ALTER TABLE `esami` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `insegna`
--

DROP TABLE IF EXISTS `insegna`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `insegna` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `id_insegnante` int DEFAULT NULL,
  `id_corso` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `ID_idx` (`id_insegnante`),
  KEY `ID_idx1` (`id_corso`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `insegna`
--

LOCK TABLES `insegna` WRITE;
/*!40000 ALTER TABLE `insegna` DISABLE KEYS */;
INSERT INTO `insegna` VALUES (1,2,1),(2,2,2),(3,1,1),(4,1,3),(5,2,0),(6,2,5),(7,11,5),(8,12,2);
/*!40000 ALTER TABLE `insegna` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `insegnante`
--

DROP TABLE IF EXISTS `insegnante`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `insegnante` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `nome` varchar(45) NOT NULL,
  `cognome` varchar(45) NOT NULL,
  `matricola` int NOT NULL,
  `password` char(64) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `insegnante`
--

LOCK TABLES `insegnante` WRITE;
/*!40000 ALTER TABLE `insegnante` DISABLE KEYS */;
INSERT INTO `insegnante` VALUES (2,'mora','morella',2,'faf4cec4bc508b56f1caec91198c3272a833310be34f481f21f471b40bf23b43'),(9,'vano','vanella',1,'d8d19ae380bd2a77c92d050d8b63a0a94b899af7a32bef617eb3e04f07a7d6f6'),(10,'test','test',3,'9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08'),(11,'Paolo','Cannone',4,'688787d8ff144c502c7f5cffaafe2cc588d86079f9de88304c26b0cb99ce91c6'),(12,'Claudio','Perrini',54,'fb78e080ef1f8b1c649f0401755d85c7770eceeb99525a8d9672c9611da37f7e');
/*!40000 ALTER TABLE `insegnante` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `risultati_esami`
--

DROP TABLE IF EXISTS `risultati_esami`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `risultati_esami` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `nome_studente` varchar(45) NOT NULL,
  `cognome_studente` varchar(45) NOT NULL,
  `contenuto` mediumtext NOT NULL,
  `ID_esame` int NOT NULL,
  `tentativi` int NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `ID_esame_idx` (`ID_esame`),
  CONSTRAINT `ID_esame` FOREIGN KEY (`ID_esame`) REFERENCES `esami` (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `risultati_esami`
--

LOCK TABLES `risultati_esami` WRITE;
/*!40000 ALTER TABLE `risultati_esami` DISABLE KEYS */;
INSERT INTO `risultati_esami` VALUES (9,'vano','vano','[\"Which anime heavily features music from the genre &quot;Eurobeat&quot;?;Initial D;Initial D\\n\",\"What is the name of the stuffed lion in Bleach?;Kon;Kon\\n\",\"What is the age of Ash Ketchum in Pokemon when he starts his journey?;10;10\\n\",\"What was Ash Ketchum&#039;s second Pokemon?;;Caterpie\\n\",\"In &quot;A Certain Scientific Railgun&quot;, how many &quot;sisters&quot; did Accelerator have to kill to achieve the rumored level 6?;128;20,000\\n\",\"On what medium was &quot;Clannad&quot; first created?;Visual novel;Visual novel\\n\",\"The characters of &quot;Log Horizon&quot; are trapped in what game?;Elder Tale;Elder Tale\\n\",\"Satella in &quot;Re:Zero&quot; is the witch of what?;Pride;Envy\\n\",\"How many &quot;JoJos&quot; that are protagonists are there in the series &quot;Jojo&#039;s Bizarre Adventure&quot;?;5;8\\n\",\"What name is the main character Chihiro given in the 2001 movie &quot;Spirited Away&quot;?;Sen (Thousand);Sen (Thousand)\\n\"]',5,1),(10,'mora','marella','[\"In &quot;Fairy Tail&quot;, what is the nickname of Natsu Dragneel?;The Salamander;The Salamander\",\"Who is the true moon princess in Sailor Moon?;Sailor Moon;Sailor Moon\",\"In &quot;Future Diary&quot;, what is the name of Yuno Gasai&#039;s Phone Diary?;Murder Diary;Yukiteru Diary\",\"In the anime Black Butler, who is betrothed to be married to Ciel Phantomhive?;Rachel Phantomhive;Elizabeth Midford\",\"In &quot;The Melancholy of Haruhi Suzumiya&quot; series, the SOS Brigade club leader is unknowingly treated as a(n) __ by her peers.;Time Traveler;God\",\"What caused the titular mascot of Yo-Kai Watch, Jibanyan, to become a yokai?;When he put on the harmaki;Being run over by a truck\",\"In the 2012 animated film &quot;Wolf Children&quot;, what are the names of the wolf children?;Ame & Yuki;Ame &amp; Yuki\",\"In Ms. Kobayashi&#039;s Dragon Maid, who is Kobayashi&#039;s maid?;Tohru;Tohru\",\"How many &quot;JoJos&quot; that are protagonists are there in the series &quot;Jojo&#039;s Bizarre Adventure&quot;?;6;8\",\"What is the theme song of &quot;Neon Genesis Evangelion&quot;?;Stardust Crusaders;A Cruel Angel&#039;s Thesis\"]',5,1),(11,'teo','tartanella','[\"In &quot;Highschool of the Dead&quot;, where did Komuro and Saeko establish to meet after the bus explosion?;Eastern Police Station;Eastern Police Station\",\"Which of the following countries does &quot;JoJo&#039;s Bizarre Adventure: Stardust Crusaders&quot; not take place in?;Philippines;Philippines\",\"What is the last line muttered in the anime film &quot;The End of Evangelion&quot;?;\\\"How disgusting.\\\";&quot;How disgusting.&quot;\",\"Winch of these names are not a character of JoJo&#039;s Bizarre Adventure?;JoJo Kikasu;JoJo Kikasu\",\"Which of these anime have over 7,500 episodes?;Sazae-san;Sazae-san\",\"Which song was the callsign for Stefan Verdemann&#039;s KWFM radio station in Urasawa Naoki&#039;s &quot;Monster&quot;?;Over the Rainbow;Over the Rainbow\",\"Medaka Kurokami from &quot;Medaka Box&quot; has what abnormality?;The End;The End\",\"In &quot;One Piece&quot;, what does &quot;the Pirate King&quot; mean to the captain of the Straw Hat Pirates?;Freedom;Freedom\",\"In &quot;One Piece&quot;, who confirms the existence of the legendary treasure, One Piece?;Edward \\\"Whitebeard\\\" Newgate;Edward &quot;Whitebeard&quot; Newgate\",\"Which one of these characters is from &quot;Legendz : Tale of the Dragon Kings&quot;?;Shiron;Shiron\",\"In &quot;Hunter x Hunter&quot;, which of the following is NOT a type of Nen aura?;Restoration;Restoration\",\"The &quot;To Love-Ru&quot; Manga was started in what year?;2006;2006\",\"Which animation studio animated &quot;Hidamari Sketch&quot;?;Shaft;Shaft\",\"Which animation studio animated &quot;To Love-Ru&quot;?;Xebec;Xebec\",\"Which animation studio animated &quot;Psycho Pass&quot;?;Production I.G;Production I.G\",\"Which animation studio produced the anime adaptation of &quot;xxxHolic&quot;?;Production I.G;Production I.G\",\"Who was the Author of the manga Uzumaki?;Junji Ito;Junji Ito\",\"Who was the Author of the manga Monster Hunter Orage?;\\tHiro Mashima;\\tHiro Mashima\",\"Who is the horror manga artist who made Uzumaki?;Junji Ito;Junji Ito\",\"Who was the Director of the 1988 Anime film &quot;Grave of the Fireflies&quot;?;Isao Takahata;Isao Takahata\"]',6,1),(12,'Davide','Vano','[\"What is the plane curve proposed by Descartes to challenge Fermat&#039;s extremum-finding techniques called?;Folium of Descartes;Folium of Descartes\",\"The French mathematician &Eacute;variste Galois is primarily known for his work in which?;Galois\' Method for PDE\'s ;Galois Theory\"]',7,1);
/*!40000 ALTER TABLE `risultati_esami` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-01-25  1:15:38