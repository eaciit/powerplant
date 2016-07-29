USE [ecsecnew]
GO
/****** Object:  StoredProcedure [dbo].[SaveSelectedPlant]    Script Date: 7/14/2016 11:57:15 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[SaveSelectedPlant]
	@Plant NVARCHAR(50) = NULL
AS
BEGIN
	SET NOCOUNT ON;

	INSERT INTO ScenarioSimulationSelectedPlant(SSID, Plant) VALUES (
	(SELECT IDENT_CURRENT ('ScenarioSimulation')),
	@Plant
)
END
