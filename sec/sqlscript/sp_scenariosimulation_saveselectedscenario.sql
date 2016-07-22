USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[SaveSelectedScenario]    Script Date: 7/14/2016 11:51:45 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[SaveSelectedScenario]
	@ID NVARCHAR(50) = NULL,
	@NAME NVARCHAR(50) = NULL,
	@VALUE FLOAT = 0
AS
BEGIN
	SET NOCOUNT ON;

	INSERT INTO ScenarioSimulationSelectedScenario
	VALUES(
		(SELECT IDENT_CURRENT ('ScenarioSimulation')),
		@ID,
		@NAME,
		@VALUE
	)
END
