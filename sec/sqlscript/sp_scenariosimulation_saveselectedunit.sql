USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[SaveSelectedUnit]    Script Date: 7/14/2016 11:58:04 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[SaveSelectedUnit] 
	@Unit NVARCHAR(50) = NULL
AS
BEGIN
	SET NOCOUNT ON;

    INSERT INTO ScenarioSimulationSelectedUnit VALUES (
	(SELECT IDENT_CURRENT ('ScenarioSimulation')),
	@Unit
)
END
