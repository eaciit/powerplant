USE [ecsecnew]
GO
/****** Object:  StoredProcedure [dbo].[SaveSelectedScenario]    Script Date: 7/14/2016 11:51:45 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[SaveSelectedScenario]
	@ID NVARCHAR(50) = NULL,
	@NAME NVARCHAR(50) = NULL,
	@VALUE FLOAT = 0
AS
BEGIN
	SET NOCOUNT ON;

	INSERT INTO ScenarioSimulationSelectedScenario (SSID, ID, NAME, VALUE)
	VALUES(
		(SELECT IDENT_CURRENT ('ScenarioSimulation')),
		@ID,
		@NAME,
		@VALUE
	)
END