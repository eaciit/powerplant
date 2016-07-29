USE [ecsecnew]
GO
/****** Object:  StoredProcedure [dbo].[GetAllScenarioSimulation]    Script Date: 7/14/2016 11:48:56 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[GetAllScenarioSimulation]
AS
BEGIN
	SET NOCOUNT ON;

    SELECT A.*, B.Plant AS 'SelectedPlant', C.Unit AS 'SelectedUnit', D.ID AS 'ScenarioID', D.Name AS 'ScenarioName', D.Value AS 'ScenarioValue' 
	FROM ScenarioSimulation AS A
	LEFT JOIN 
	ScenarioSimulationSelectedPlant AS B
	ON A.Id = B.SSId
	LEFT JOIN
	ScenarioSimulationSelectedUnit AS C
	ON A.Id = C.SSId
	LEFT JOIN
	ScenarioSimulationSelectedScenario AS D
	ON A.Id = D.SSId
END
