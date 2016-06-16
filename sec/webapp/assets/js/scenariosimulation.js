// Set preset value
ScenarioSimulation.HistoricValueEquation = {
    Revenue:ko.observable(0),
    LaborCost:ko.observable(0),
    MaterialCost:ko.observable(0),
    ServiceCost:ko.observable(0),
    OperatingCost:ko.observable(0),
};
ScenarioSimulation.FutureValueEquation = {
    Revenue:ko.observable(0),
    LaborCost:ko.observable(0),
    MaterialCost:ko.observable(0),
    ServiceCost:ko.observable(0),
    OperatingCost:ko.observable(0),
};
ScenarioSimulation.Differential = {
};
ScenarioSimulation.HistoricValueEquation.MaintenanceCost = ko.computed(function(){
    var d = ScenarioSimulation.HistoricValueEquation;
    return d.LaborCost()+d.MaterialCost()+d.ServiceCost();
});
ScenarioSimulation.HistoricValueEquation.ValueEquation = ko.computed(function(){
    var d = ScenarioSimulation.HistoricValueEquation;
    return d.Revenue()-d.MaintenanceCost()-d.OperatingCost();
});

ScenarioSimulation.FutureValueEquation.MaintenanceCost = ko.computed(function(){
    var d = ScenarioSimulation.FutureValueEquation;
    return d.LaborCost()+d.MaterialCost()+d.ServiceCost();
});
ScenarioSimulation.FutureValueEquation.ValueEquation = ko.computed(function(){
    var d = ScenarioSimulation.FutureValueEquation;
    return d.Revenue()-d.MaintenanceCost()-d.OperatingCost();
});

// Differential
ScenarioSimulation.Differential.Revenue = ko.computed(function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    return h.Revenue() == 0 ? 0 : (f.Revenue()-h.Revenue())/h.Revenue();
});
ScenarioSimulation.Differential.LaborCost = ko.computed(function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    return h.LaborCost() == 0 ? 0 : (f.LaborCost()-h.LaborCost())/h.LaborCost();
});
ScenarioSimulation.Differential.MaterialCost = ko.computed(function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    return h.MaterialCost() == 0 ? 0 : (f.MaterialCost()-h.MaterialCost())/h.MaterialCost();
});
ScenarioSimulation.Differential.ServiceCost = ko.computed(function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    return h.ServiceCost() == 0 ? 0 : (f.ServiceCost()-h.ServiceCost())/h.ServiceCost();
});

ScenarioSimulation.Differential.MaintenanceCost = ko.computed(function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    return h.MaintenanceCost() == 0 ? 0 : (f.MaintenanceCost()-h.MaintenanceCost())/h.MaintenanceCost();
});
ScenarioSimulation.Differential.OperatingCost = ko.computed(function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    return h.OperatingCost() == 0 ? 0 : (f.OperatingCost()-h.OperatingCost())/h.OperatingCost();
});
ScenarioSimulation.Differential.ValueEquation = ko.computed(function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    return h.ValueEquation() == 0 ? 0 : (f.ValueEquation()-h.ValueEquation())/h.ValueEquation();
});

ScenarioSimulation.ResetData = function(){
    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    h.Revenue(0);
    h.LaborCost(0);
    h.MaterialCost(0);
    h.ServiceCost(0);
    h.OperatingCost(0);

    f.Revenue(0);
    f.LaborCost(0);
    f.MaterialCost(0);
    f.ServiceCost(0);
    f.OperatingCost(0);
};
ScenarioSimulation.PrintSimulation = function(){
    ScenarioSimulation.isPrinting(true);
    selector = "#ValueEquationSimulation";
    fname = "Value Equation Simulation.pdf"
    // Convert the DOM element to a drawing using kendo.drawing.drawDOM
    kendo.drawing.drawDOM($(selector))
    .then(function (group) {
        // Render the result as a PDF file
        return kendo.drawing.exportPDF(group, {
            paperSize: "auto",
            margin: { left: "1cm", top: "1cm", right: "1cm", bottom: "1cm" }
        });
    })
    .done(function (data) {
        // Save the PDF file
        kendo.saveAs({
            dataURI: data,
            fileName: fname,
        });
    });
    ScenarioSimulation.isPrinting(false);
}
ScenarioSimulation.SaveSimulation = function(){
    if(!confirm("Are you sure ?")){
        return false;
    }
    var url = "/scenariosimulation/savedata"
    var parm = ScenarioSimulation.GetFilter();
    parm.SimulationName = ScenarioSimulation.SimulationName();
    parm.SimulationDescription = ScenarioSimulation.SimulationDescription();

    var h = ScenarioSimulation.HistoricValueEquation;
    var f = ScenarioSimulation.FutureValueEquation;
    var d = ScenarioSimulation.Differential;

    var ScenarioList = ScenarioSimulation.scenarioList();
    var SelectedScenario =  Enumerable.From(ScenarioList).Where(function(x){return x.isSelected()}).ToArray();
    parm.SelectedScenario = [];
    for(var i in SelectedScenario){
        parm.SelectedScenario.push({ID:SelectedScenario[i].ID,Name:SelectedScenario[i].Name,Value:SelectedScenario[i].Value()})
    }
    parm.SelectedScenarioLength = SelectedScenario.length;
    parm.HistoricData = {
        Revenue:h.Revenue(),
        LaborCost:h.LaborCost(),
        MaterialCost:h.MaterialCost(),
        ServiceCost:h.ServiceCost(),
        OperatingCost:h.OperatingCost(),
        MaintenanceCost:h.MaintenanceCost(),
        ValueEquation:h.ValueEquation(),
    };
    parm.FutureData = {
        Revenue:f.Revenue(),
        LaborCost:f.LaborCost(),
        MaterialCost:f.MaterialCost(),
        ServiceCost:f.ServiceCost(),
        OperatingCost:f.OperatingCost(),
        MaintenanceCost:f.MaintenanceCost(),
        ValueEquation:f.ValueEquation(),
    };

    parm.Differential = {
        Revenue:d.Revenue(),
        LaborCost:d.LaborCost(),
        MaterialCost:d.MaterialCost(),
        ServiceCost:d.ServiceCost(),
        OperatingCost:d.OperatingCost(),
        MaintenanceCost:d.MaintenanceCost(),
        ValueEquation:d.ValueEquation(),
    };
    ScenarioSimulation.Processing(true);
   ajaxPost(url, parm, function(data){
    //console.log(data);
      if(data.Status=="OK"){   
            ScenarioSimulation.GetDataSimulation();
            ScenarioSimulation.Processing(false);
            ScenarioSimulation.ResetSelection();
            ScenarioSimulation.isCreatingSimulation(false);
      }else{
          alert(data.Message);
      }
  });
}
ScenarioSimulation.ProcessingData = function(dataSource){
    var FutureDataSource = dataSource;
    //console.log(dataSource);
    // Pre Set Value
    var HistoricData = ScenarioSimulation.HistoricValueEquation;
    var FutureData = ScenarioSimulation.FutureValueEquation;
    if (dataSource!== undefined){
        HistoricData.Revenue(Enumerable.From(dataSource).Sum(function(x){return x.Revenue}));
        HistoricData.LaborCost(Enumerable.From(dataSource).Sum(function(x){return x.TotalLabourCost}));
        HistoricData.MaterialCost(Enumerable.From(dataSource).Sum(function(x){return x.TotalMaterialCost}));
        HistoricData.ServiceCost(Enumerable.From(dataSource).Sum(function(x){return x.TotalServicesCost}));
        HistoricData.OperatingCost(Enumerable.From(dataSource).Sum(function(x){return x.OperatingCost}));
    }
    
    var ScenarioList = ScenarioSimulation.scenarioList();
    var SelectedScenario =  Enumerable.From(ScenarioList).Where(function(x){return x.isSelected()}).ToArray();
    // console.log(SelectedScenario);

    // Calculating Future Value
    for(var i in SelectedScenario){
        switch(SelectedScenario[i].ID){
            case "ReduceMaterialCost" :
                ScenarioSimulation.ReduceMaterialCost(SelectedScenario[i].Value(),FutureDataSource);
                break;
            case "ReduceOutages":
                ScenarioSimulation.ReduceOutages(SelectedScenario[i].Value(),FutureDataSource);
                break;
            case "ReduceMaintenanceDuration":
                ScenarioSimulation.ReduceMaintenanceDuration(SelectedScenario[i].Value(),FutureDataSource);
                break;
            default:break;
        }
    }
    if (FutureDataSource!== undefined){
        FutureData.Revenue(Enumerable.From(FutureDataSource).Sum(function(x){return x.Revenue}));
        FutureData.LaborCost(Enumerable.From(FutureDataSource).Sum(function(x){return x.TotalLabourCost}));
        FutureData.MaterialCost(Enumerable.From(FutureDataSource).Sum(function(x){return x.TotalMaterialCost}));
        FutureData.ServiceCost(Enumerable.From(FutureDataSource).Sum(function(x){return x.TotalServicesCost}));
        FutureData.OperatingCost(Enumerable.From(FutureDataSource).Sum(function(x){return x.OperatingCost}));
    }
}
ScenarioSimulation.ReduceMaterialCost = function(value,FutureDataSource){
    for(var i in FutureDataSource){
        FutureDataSource[i].TotalMaterialCost = (100-value)/100*FutureDataSource[i].TotalMaterialCost;
    }
}
ScenarioSimulation.ReduceOutages = function(value,FutureDataSource){
    for(var i in FutureDataSource){
        var Increases = value/100*FutureDataSource[i].AvgNetGeneration*(FutureDataSource[i].TotalOutageDuration/24)*FutureDataSource[i].VOMR;
        FutureDataSource[i].Revenue = FutureDataSource[i].Revenue+Increases;
    }
}
ScenarioSimulation.ReduceMaintenanceDuration = function(value,FutureDataSource){
    for(var i in FutureDataSource){
        FutureDataSource[i].TotalLabourCost = (100-value)/100*FutureDataSource[i].TotalLabourCost;
        FutureDataSource[i].TotalMaterialCost = (100-value)/100*FutureDataSource[i].TotalMaterialCost;
        FutureDataSource[i].TotalServicesCost = (100-value)/100*FutureDataSource[i].TotalServicesCost;
        var Increases = value/100*FutureDataSource[i].AvgNetGeneration*(FutureDataSource[i].TotalOutageDuration/24)*FutureDataSource[i].VOMR;
        FutureDataSource[i].Revenue = FutureDataSource[i].Revenue+Increases;

        var OperatingIncreasing = (FutureDataSource[i].TotalFuelCost / (365*24)) * ((value/100) * FutureDataSource[i].TotalOutageDuration)
        FutureDataSource[i].OperatingCost += OperatingIncreasing 
    }
}
ScenarioSimulation.SimulationData = ko.observableArray([]);
ScenarioSimulation.GenDataSimulation = function(dataSource){
    for(var i in dataSource){
        dataSource[i].Index = i;
        dataSource[i].SelectedPlantStr = "";
        dataSource[i].SelectedUnitStr = "";
        for (var x in dataSource[i].SelectedPlant){
            if(x == dataSource[i].SelectedPlant.length-1){
                dataSource[i].SelectedPlantStr += dataSource[i].SelectedPlant[x];
            }else{
                dataSource[i].SelectedPlantStr += dataSource[i].SelectedPlant[x] + " , ";
            }
        }
        for (var x in dataSource[i].SelectedUnit){
            if(x == dataSource[i].SelectedUnit.length-1){
                dataSource[i].SelectedUnitStr += dataSource[i].SelectedUnit[x];
            }else{
                dataSource[i].SelectedUnitStr += dataSource[i].SelectedUnit[x] + " , ";
            }
        }
    }
    ScenarioSimulation.SimulationData(dataSource);
    $("#SimulationData").html("");
    $("#SimulationData").kendoGrid({
        dataSource: {
            data: dataSource,
            pageSize: 20
        },
        scrollable: true,
        sortable: true,
        filterable: true,
        columns: [
            { field:"Name",title: "Simulation Name", width:250},
            { field:"Description",title: "Simulation Description"},
            { title: "Period", template:"#:kendo.toString(getUTCDate(Start_Period),'dd MMM yyyy')# - #:kendo.toString(getUTCDate(End_Period),'dd MMM yyyy')#",width:150},
            { title: "Selected Plant",template:"#:SelectedPlantStr#",width:200},
            { title: "Selected Unit",template:"#:SelectedUnitStr#",width:200},
            { title: "",template:"<button type='button' class='btn btn-default btn-xs' onclick='ScenarioSimulation.GetDetailSimulation(#:Index#)'><span class='fa fa-newspaper-o'></span></button>&nbsp;<button type='button' class='btn btn-default btn-xs remove-button' onclick='ScenarioSimulation.RemoveSimulation(#:Index#)'><span class='fa fa-times'></span></button>",width:70},
        ]
    });
}
ScenarioSimulation.RemoveSimulation = function(index){
    var data = ScenarioSimulation.SimulationData()[index];
    ScenarioSimulation.selectedSimulation(data.Name);
    ScenarioSimulation.selectedDescription(data.Description);
    if(!confirm("Are you sure ?")){
        return false;
    }
    var url = "/scenariosimulation/removedata"
    var parm = ScenarioSimulation.GetFilter();
    ajaxPost(url, parm, function(data){
      if(data.Status=="OK"){   
            ScenarioSimulation.selectedSimulation("");
            ScenarioSimulation.selectedDescription("");
            ScenarioSimulation.GetDataSimulation();
      }else{
          alert(data.Message);
      }
  });
}
ScenarioSimulation.GetDetailSimulation = function(index){
    ScenarioSimulation.isCreatingSimulation(true);
    ScenarioSimulation.isResultAvailable(true);
    ScenarioSimulation.isSelectedAvailable(true);
    var data = ScenarioSimulation.SimulationData()[index];
    ScenarioSimulation.selectedSimulation(data.Name);
    ScenarioSimulation.selectedDescription(data.Description);
    ScenarioSimulation.SimulationName(data.Name);
    ScenarioSimulation.SimulationDescription(data.Description);
    var scenarioList = ScenarioSimulation.scenarioList();
    for(var i in data.SelectedScenario){
        for(var j in scenarioList){
            if(data.SelectedScenario[i].ID == scenarioList[j].ID){
                scenarioList[j].isSelected(true);
                scenarioList[j].isSet(true);
                scenarioList[j].Value(data.SelectedScenario[i].Value);
            }
        }   
    }
    ScenarioSimulation.RunSimulation();
}