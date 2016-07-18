USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[SaveSummaryData]    Script Date: 7/14/2016 11:58:45 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[SaveSummaryData]
	@FunctionalLocation NVARCHAR(50) = NULL,
	@FLDescription TEXT = NULL,
	@SortField NVARCHAR(50) = NULL,
	@ParentFL NVARCHAR(50) = NULL,
	@HasChild BIT = 0,
	@Province NVARCHAR(50) = NULL,
	@Region NVARCHAR(50) = NULL,
	@City NVARCHAR(50) = NULL,
	@GasTurbineUnit INT = 0,
	@GasTurbineCapacity FLOAT = 0,
	@SteamUnit INT = 0,
	@SteamUnitCapacity FLOAT = 0,
	@DieselUnit INT = 0,
	@DieselUnitCapacity FLOAT = 0,
	@CombinedCycleUnit INT = 0,
	@CombinedCycleUnitCapacity FLOAT = 0
AS
BEGIN
	SET NOCOUNT ON;

    INSERT INTO SummaryData
	VALUES(
		@FunctionalLocation,
		@FLDescription,
		@SortField,
		@ParentFL,
		@HasChild,
		@Province,
		@Region,
		@City,
		@GasTurbineUnit,
		@GasTurbineCapacity,
		@SteamUnit,
		@SteamUnitCapacity,
		@DieselUnit,
		@DieselUnitCapacity,
		@CombinedCycleUnit,
		@CombinedCycleUnitCapacity
	)
END
