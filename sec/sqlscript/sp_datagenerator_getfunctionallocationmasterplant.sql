USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[GetFunctionalLocationMasterPlant]    Script Date: 7/14/2016 11:54:48 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[GetFunctionalLocationMasterPlant]
AS
BEGIN
	SET NOCOUNT ON;

    SELECT * FROM FunctionalLocation WHERE LEN(FunctionalLocationCode) = 4
END
