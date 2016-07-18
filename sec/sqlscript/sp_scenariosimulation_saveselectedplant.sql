USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[SaveSelectedPlant]    Script Date: 7/14/2016 11:57:15 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[SaveSelectedPlant]
	@Plant NVARCHAR(50) = NULL
AS
BEGIN
	SET NOCOUNT ON;

	INSERT INTO ScenarioSimulationSelectedPlant VALUES (
	(SELECT IDENT_CURRENT ('ScenarioSimulation')),
	@Plant
)
END
