USE [ecsec]
GO
/****** Object:  StoredProcedure [dbo].[DataBrowserSP]    Script Date: 6/9/2016 2:53:02 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Faisal Reza
-- Create date: 9 June 2016
-- Description:	To get Data Browser data for H3
-- =============================================
CREATE PROCEDURE [dbo].[SP_DataBrowser_H3] 
	@ORDER nvarchar(100) = 'PlantPlantName',
	@DIR nvarchar(10) = 'asc',
	@Offset int = 1,
	@Limit int = 10,
	@PeriodYear int = null,
	@PeriodFrom nvarchar(10) = null,
	@PeriodTo nvarchar(10) = null,
	@EquipmentType nvarchar(50) = null,
	@EquipmentTypeEmpty nvarchar(50) = null,
	@PlantName nvarchar(30) = null,
	@PlantNameEmpty nvarchar(30) = null,
	@PlantCode nvarchar(30) = null,
	@PlantCodeEmpty nvarchar(30) = null,
	@MaintenanceFunctionalLocation nvarchar(50) = null,
	@WorkOrderType nvarchar(10) = null,	
	@FILTERS nvarchar(2000) = null	
AS
BEGIN
	IF @EquipmentType IS NULL
		BEGIN
			SET @EquipmentTypeEmpty = 'xxx'	
		END

	IF @PlantName IS NULL
		BEGIN
			SET @PlantNameEmpty = ''
		END

	SET NOCOUNT ON;

	DECLARE @SQLQuery nvarchar(max)
    DECLARE @ParamDefinition nvarchar(1000)

	SET @ParamDefinition = ' @ORDER nvarchar(100),
		@DIR nvarchar(10),
		@Offset int,
		@Limit int = 10,
		@PeriodYear int = null,
		@PeriodFrom nvarchar(10) = null,
		@PeriodTo nvarchar(10) = null,
		@EquipmentType nvarchar(50) = null,
		@EquipmentTypeEmpty nvarchar(50) = null,
		@PlantName nvarchar(30) = null,
		@PlantNameEmpty nvarchar(30) = null,
		@PlantCode nvarchar(30) = null,
		@PlantCodeEmpty nvarchar(30) = null,
		@MaintenanceFunctionalLocation nvarchar(50) = null,
		@WorkOrderType nvarchar(10) = null,
		@FILTERS nvarchar(2000) = null'

	SET @SQLQuery = '
				SELECT
					*,
					(Select CAST((Select case
							when res.MaintenanceIntervalTmp > 0 or res.MaintenanceIntervalTmp is not null
							then res.MaintenanceIntervalTmp
							else 0
						END) as float)
					) as MaintenanceInterval
				from(
					Select
						DataBrowser.Id as Id,
						DataBrowser.PeriodYear as PeriodYear,
						DataBrowser.FunctionalLocation as FunctionalLocation,
						DataBrowser.FLDescription as FLDescription,
						DataBrowser.IsTurbine as IsTurbine,
						DataBrowser.IsSystem as IsSystem,
						DataBrowser.TurbineParent as TurbineParent,
						DataBrowser.SystemParent SystemParent,
						DataBrowser.AssetType as AssetType,
						DataBrowser.EquipmentType EquipmentType,
						DataBrowser.EquipmentTypeDescription as EquipmentTypeDescription,
						DataBrowser.PlantCode as PlantCode,
						DataBrowser.TInfShortName as TInfShortName,
						DataBrowser.TInfManufacturer as TInfManufacturer,
						DataBrowser.TInfModel as TInfModel,
						DataBrowser.TInfUnitType as TInfUnitType,
						DataBrowser.TInfInstalledCapacity as TInfInstalledCapacity,
						DataBrowser.TInfOperationalCapacity as TInfOperationalCapacity,
						DataBrowser.TInfPrimaryFuel as TInfPrimaryFuel,      
						DataBrowser.TInfPrimaryFuel2 as TInfPrimaryFuel2,     
						DataBrowser.TInfBackupFuel as TInfBackupFuel,       
						DataBrowser.TInfHeatRate as TInfHeatRate,
						DataBrowser.TInfEfficiency as TInfEfficiency,       
						DataBrowser.TInfCommisioningDate as TInfCommisioningDate, 
						DataBrowser.TInfRetirementPlan as TInfRetirementPlan,   
						DataBrowser.TInfInstalledMWH as TInfInstalledMWH,     
						DataBrowser.TInfActualEnergyGeneration as TInfActualEnergyGeneration,
						DataBrowser.TInfActualFuelConsumption_GASMMSCF as TInfActualFuelConsumption_GASMMSCF,
						DataBrowser.TInfActualFuelConsumption_CrudeBarrel as TInfActualFuelConsumption_CrudeBarrel,
						DataBrowser.TInfActualFuelConsumption_HFOBarrel as TInfActualFuelConsumption_HFOBarrel,
						DataBrowser.TInfActualFuelConsumption_DieselBarrel as TInfActualFuelConsumption_DieselBarrel,
						DataBrowser.TInfCapacityFactor as TInfCapacityFactor,   
						DataBrowser.TInfUpdateEnergyGeneration as TInfUpdateEnergyGeneration,
						DataBrowser.TInfUpdateFuelConsumption as TInfUpdateFuelConsumption, 

						Plant.PlantCode as PlantPlantCode,
						Plant.PlantName as PlantPlantName,
						Plant.PlantType as PlantPlantType,
						Plant.Province as PlantProvince,
						Plant.Region as PlantRegion,
						Plant.City as PlantCity,
						Plant.FuelTypes_Crude as PlantFuelTypes_Crude,
						Plant.FuelTypes_Heavy as PlantFuelTypes_Heavy,
						Plant.FuelTypes_Diesel as PlantFuelTypes_Diesel,
						Plant.FuelTypes_Gas as PlantFuelTypes_Gas,
						Plant.GasTurbineUnit as PlantGasTurbineUnit,
						Plant.GasTurbineCapacity as PlantGasTurbineCapacity,
						Plant.SteamUnit as PlantSteamUnit,
						Plant.SteamCapacity as PlantSteamCapacity,
						Plant.DieselUnit as PlantDieselUnit,
						Plant.DieselCapacity as PlantDieselCapacity,
						Plant.CombinedCycleUnit as PlantCombinedCycleUnit,
						Plant.CombinedCycleCapacity as PlantCombinedCycleCapacity,
						Plant.Longitude as PlantLongitude,
						Plant.Latitude as PlantLatitude,
 		
						WO.Type as WorkOrderType,
						WO.FunctionalLocation as MaintenanceFunctionalLocation,	
						WO.OrderCode as MaintenanceOrder,
						WO.Description as MaintenanceDescription,
						WO.ScheduledStart as PlanStart,
						WO.ScheduledFinish as PlanFinish,
						WO.ActualStart as ActualStart,
						WO.ActualFinish as ActualFinish,
						WO.ActualCost as MaintenanceCost,
						(Select top 1 mc.MaintenanceActivityType
							from MaintenanceCost as mc
							where mc.MaintenanceOrder = WO.OrderCode) as ActivityType,				
						(Select CAST((Select case
								when DateDiff(MINUTE, WO.ScheduledStart,WO.ScheduledFinish) > 0 
								then DateDiff(MINUTE, WO.ScheduledStart,WO.ScheduledFinish) / 60.0
								else 0
							END) as float)
						) as PlanDuration,
						(Select CAST((Select case
								when DateDiff(MINUTE, WO.ActualStart, WO.ActualFinish) > 0
								then DateDiff(MINUTE, WO.ActualStart, WO.ActualFinish) / 60.0
								else 0
							END) as float)
						) as ActualDuration,
						(SELECT CAST((DATEDIFF(
							HOUR,
							(Select top 1 wl.ActualFinish from WOList as wl where wl.FunctionalLocation = DataBrowser.FunctionalLocation and wl.ActualStart < WO.ActualStart order by wl.ActualStart desc),
							WO.ActualStart) / 24.0) as float)) as MaintenanceIntervalTmp
					from DataBrowser DataBrowser
					LEFT JOIN PowerPlantCoordinates Plant 
					ON DataBrowser.PlantCode = Plant.PlantCode
					LEFT Join WOList WO
					ON DataBrowser.FunctionalLocation = WO.FunctionalLocation
					WHERE 
						PeriodYear = ISNULL(@PeriodYear,PeriodYear) and

						WO.ActualStart >= ISNULL(@PeriodFrom,WO.ActualStart) and
						WO.ActualStart < ISNULL(@PeriodTo,WO.ActualStart) and
			
						EquipmentType IN (ISNULL(@EquipmentType,EquipmentType)) and 
						EquipmentType <> ISNULL(@EquipmentTypeEmpty,'') and 
			
						PlantName IN (ISNULL(@PlantName,PlantName)) and
						PlantName <> ISNULL(@PlantNameEmpty,'') and

						Plant.PlantCode = ISNULL(@PlantCode,Plant.PlantCode) and
						WO.FunctionalLocation = ISNULL(@MaintenanceFunctionalLocation,WO.FunctionalLocation) and
						WO.Type IN (ISNULL(@WorkOrderType,WO.Type))'

	IF @ORDER IS NOT NULL
		SET @SQLQuery = @SQLQuery + ' ORDER BY DataBrowser.ID ASC '

	--IF @DIR IS NOT NULL 
	--	SET @SQLQuery = @SQLQuery + ' @DIR ' 
		
	SET @SQLQuery = @SQLQuery + ' OFFSET @Offset ROWS '
	SET @SQLQuery = @SQLQuery + ' FETCH NEXT @Limit ROWS ONLY '

	SET @SQLQuery = @SQLQuery + ') as res'
	
	IF @FILTERS IS NOT NULL
		SET @SQLQuery = @SQLQuery + ' @FILTERS '

	EXECUTE sp_Executesql @SQLQuery,
				@ParamDefinition,
				@ORDER,
				@DIR,
				@Offset,
				@Limit,
				@PeriodYear,
				@PeriodFrom,
				@PeriodTo,
				@EquipmentType,
				@EquipmentTypeEmpty,
				@PlantName,
				@PlantNameEmpty,
				@PlantCode,
				@PlantCodeEmpty,
				@MaintenanceFunctionalLocation,
				@WorkOrderType,
				@FILTERS

	
END
