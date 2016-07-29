USE [ecsecnew]
GO
/****** Object:  StoredProcedure [dbo].[DeleteScenarioSimulation]    Script Date: 7/14/2016 11:49:54 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[DeleteScenarioSimulation] 
	@NAME NVARCHAR(50) = NULL,
	@DESC NVARCHAR(50) = NULL
AS
BEGIN
	SET NOCOUNT ON;

    DELETE ScenarioSimulationSelectedPlant WHERE SSId = (SELECT Id FROM ScenarioSimulation WHERE Name=@NAME AND Description=@DESC)
	DELETE ScenarioSimulationSelectedUnit WHERE SSId = (SELECT Id FROM ScenarioSimulation WHERE Name=@NAME AND Description=@DESC)
	DELETE ScenarioSimulationSelectedScenario WHERE SSId = (SELECT Id FROM ScenarioSimulation WHERE Name=@NAME AND Description=@DESC)
	DELETE ScenarioSimulation WHERE Name=@NAME AND Description=@DESC
END
