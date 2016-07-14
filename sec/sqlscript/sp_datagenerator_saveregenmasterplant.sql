USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[SaveRegenMasterPlant]    Script Date: 7/14/2016 11:56:36 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER PROCEDURE [dbo].[SaveRegenMasterPlant]
	@PlantCode NVARCHAR(50) = NULL,
	@PlantName NVARCHAR(50) = NULL,
	@PlantType NVARCHAR(50) = NULL,
	@Province NVARCHAR(50) = NULL,
	@Region NVARCHAR(50) = NULL,
	@City NVARCHAR(50) = NULL,
	@FuelTypes_Crude BIT = 0,
	@FuelTypes_Heavy BIT = 0,
	@FuelTypes_Diesel BIT = 0,
	@FuelTypes_Gas BIT = 0,
	@GasTurbineUnit INT = 0,
	@GasTurbineCapacity FLOAT = 0,
	@SteamUnit INT = 0,
	@SteamCapacity FLOAT = 0,
	@DieselUnit INT = 0,
	@DieselCapacity FLOAT = 0,
	@CombinedCycleUnit INT = 0,
	@CombinedCycleCapacity FLOAT = 0
AS
BEGIN
	SET NOCOUNT ON;

    INSERT INTO RegenMasterPlant 
	VALUES(
		@PlantCode,
		@PlantName,
		@PlantType,
		@Province,
		@Region,
		@City,
		@FuelTypes_Crude,
		@FuelTypes_Heavy,
		@FuelTypes_Diesel,
		@FuelTypes_Gas,
		@GasTurbineUnit,
		@GasTurbineCapacity,
		@SteamUnit,
		@SteamCapacity,
		@DieselUnit,
		@DieselCapacity,
		@CombinedCycleUnit,
		@CombinedCycleCapacity
	)
END
