USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[GetFunctionalLocation]    Script Date: 7/14/2016 11:52:50 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[GetFunctionalLocation] 
AS
BEGIN
	SET NOCOUNT ON;

	SELECT TOP 100000 * FROM FunctionalLocation
END
