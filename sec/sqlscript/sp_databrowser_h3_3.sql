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
ALTER PROCEDURE [dbo].[SP_DataBrowser_H3] 
	@ORDER nvarchar(max) = 'PlantPlantName',
	@DIR nvarchar(max) = 'asc',
	@Offset int = 1,
	@Limit int = 10,
	@PeriodYear int = null,
	@PeriodFrom nvarchar(max) = null,
	@PeriodTo nvarchar(max) = null,
	@EquipmentType nvarchar(max) = null,
	@EquipmentTypeEmpty nvarchar(max) = null,
	@PlantName nvarchar(max) = null,
	@PlantNameEmpty nvarchar(max) = null,
	@PlantCode nvarchar(max) = null,
	@PlantCodeEmpty nvarchar(max) = null,
	@WorkOrderType nvarchar(max) = null,	
	@FILTERS_DB nvarchar(max) = null,
    @FILTERS_WO nvarchar(max) = null,	
    @FILTERS_PLANT nvarchar(max) = null,
	@FILTERS_MAIN nvarchar(max) = null
AS
BEGIN
	IF (@EquipmentType IS NULL)
			SET @EquipmentTypeEmpty = 'xxx'

	IF (@PlantName IS NULL)
			SET @PlantNameEmpty = ''

	SET NOCOUNT ON;

	DECLARE @SQLQuery nvarchar(max)
    DECLARE @ParamDefinition nvarchar(max)

	SET @ParamDefinition = ' @ORDER nvarchar(max),
		@DIR nvarchar(max),
		@Offset int,
		@Limit int = 10,
		@PeriodYear int = null,
		@PeriodFrom nvarchar(max) = null,
		@PeriodTo nvarchar(max) = null,
		@EquipmentType nvarchar(max) = null,
		@EquipmentTypeEmpty nvarchar(max) = null,
		@PlantName nvarchar(max) = null,
		@PlantNameEmpty nvarchar(max) = null,
		@PlantCode nvarchar(max) = null,
		@PlantCodeEmpty nvarchar(max) = null,
		@WorkOrderType nvarchar(max) = null,
		@FILTERS_DB nvarchar(max) = null,
    	@FILTERS_WO nvarchar(max) = null,	
    	@FILTERS_PLANT nvarchar(max) = null,
		@FILTERS_MAIN nvarchar(max) = null
	'

	SET @SQLQuery = '
	;WITH CTE_DB as 
	(
		SELECT * FROM DataBrowser
		WHERE
			PeriodYear = ISNULL(@PeriodYear,PeriodYear) and
			PlantCode = ISNULL(@PlantCode,PlantCode) and
			EquipmentType IN (ISNULL(@EquipmentType,EquipmentType)) and 
			EquipmentType <> ISNULL(@EquipmentTypeEmpty,'')
	'

	IF (@FILTERS_DB IS NOT NULL)
		SET @SQLQuery = @SQLQuery + ' @FILTERS_DB '

	SET @SQLQuery = @SQLQuery + ')'

	SET @SQLQuery = @SQLQuery + '
	,
	CTE_PLANT as
	(
		SELECT * FROM PowerPlantCoordinates
		WHERE 
			PlantName IN (ISNULL(@PlantName,PlantName)) and
			PlantName <> ISNULL(@PlantNameEmpty,'') and
			PlantCode = ISNULL(@PlantCode,PlantCode)
	'

	IF (@FILTERS_PLANT IS NOT NULL)
		SET @SQLQuery = @SQLQuery + ' @FILTERS_PLANT '

	SET @SQLQuery = @SQLQuery + '),'

	SET @SQLQuery = @SQLQuery + '
	CTE_WO as 
	(
		SELECT * from 
		(SELECT *,
			/*(Select top 1 mc.MaintenanceActivityType
				from MaintenanceCost as mc
				where mc.MaintenanceOrder = WO.OrderCode) as ActivityType,*/
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
			) as ActualDuration
		FROM WOList as WO
		WHERE 
			WO.ActualStart >= ISNULL(@PeriodFrom,WO.ActualStart) and
			WO.ActualStart < ISNULL(@PeriodTo,WO.ActualStart) and
			WO.Type IN (ISNULL(@WorkOrderType,WO.Type))
		) as res
		WHERE 1=1
	'

	IF (@FILTERS_WO IS NOT NULL)
		SET @SQLQuery = @SQLQuery + ' @FILTERS_WO '

	SET @SQLQuery = @SQLQuery + ') '	

	SET @SQLQuery = @SQLQuery + '
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

				WO.PlanDuration as PlanDuration,
				WO.ActualDuration as ActualDuration,
				(Select case
						when WO.ActualStart IS NOT NULL
						then (SELECT CAST((DATEDIFF(
								HOUR,
								(Select top 1 wl.ActualFinish from WOList as wl where wl.FunctionalLocation = DataBrowser.FunctionalLocation and wl.ActualStart < WO.ActualStart order by wl.ActualStart desc),
								WO.ActualStart) / 24.0) as float))
						else 0
					END) as MaintenanceInterval						
							
			from CTE_DB as DataBrowser 
			LEFT JOIN CTE_PLANT as Plant ON DataBrowser.PlantCode = Plant.PlantCode 
			LEFT Join CTE_WO as WO ON DataBrowser.FunctionalLocation = WO.FunctionalLocation
			WHERE 
				DataBrowser.FunctionalLocation = WO.FunctionalLocation and
				DataBrowser.PlantCode = Plant.PlantCode
	'

	IF (@FILTERS_MAIN IS NOT NULL)
		SET @SQLQuery = @SQLQuery + ' @FILTERS_MAIN '

	SET @SQLQuery = @SQLQuery + ' ORDER BY Plant.PlantName ASC '	
	SET @SQLQuery = @SQLQuery + ' OFFSET @Offset ROWS '
	SET @SQLQuery = @SQLQuery + ' FETCH NEXT @Limit ROWS ONLY '

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
				@WorkOrderType,
				@FILTERS_DB,
    			@FILTERS_WO,	
    			@FILTERS_PLANT,
				@FILTERS_MAIN

	
END