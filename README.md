# Demo of Azure integration with Event Hubs, Azure Stream Analytics, and Azure SQL

## Problem Statement

A customer is sending messages of different types to Event Hubs and wants to save these messages to Azure SQL tables without having to write custom code.

## Solution with Stream Analytics

Azure Stream Analytics offers functionaltiy that allows messages to be read from Event Hubs, filter them, process them and save them to different destinations including Azure SQL tables.

## Solution Diagram

![Solution Diagram](images/Architecture-Hubs-Stream-Sql.png)

## Services Deployment

[Bicep template](deployment/main.bicep)

## Services

### Events Generator

For the purposes of this demo, the generator code has been implemented using Azure Functions with .NET Core 3.

#### Event Classes

```c#
public class BaseEvent
{
    public BaseEvent()
    {
        Ts = DateTime.UtcNow;
    }

    public string DeviceId { get; set; }
    public string Type { get; set; }
    public DateTime Ts { get; set; }
}

public class ACEvent : BaseEvent
{
    public ACMessage()
    {
        Type = typeof(ACMessage).Name;
    }
    public float CoolantTemperature { get; set; }
    public float AirFlow { get; set; }
    public float AirTemperature { get; set; }
}

public class GeneratorEvent : BaseEvent
{
    public GeneratorEvent()
    {
        Type = typeof(GeneratorEvent).Name;
    }
    public float Hertz { get; set; }
    public float Amps { get; set; }
    public float Voltage { get; set; }
    public float GasPercentage { get; set; }

}

public class MotorEvent : BaseEvent
{
    public MotorEvent()
    {
        Type = typeof(MotorEvent).Name;            
    }
    public float Temperature { get; set; }
    public float Revolutions { get; set; }
}
```

#### Emmiting Events to Event Hubs

```go
func GetRandomEvent() string {
	anomaly := getRandom(1, 6)
	eventType := getRandom(1, 4)

	airTemperature := float64(getRandom(600, 800)) / 10.0
	airFlow := float64(getRandom(30, 40)) / 10.0
	coolantTemperature := float64(getRandom(200, 400)) / 10.0
	gasPercentage := float64(getRandom(1, 1000)) / 10.0
	voltage := float64(getRandom(2300, 2450)) / 10.0
	motorTemp := float64(getRandom(1800, 2000)) / 10.0
	motorRevolutions := float64(getRandom(2000, 5000)) / 10.0
	hertz := float64(getRandom(580, 650)) / 10.0
	amps := float64(getRandom(150, 250)) / 10.0

	if anomaly == RaiseAnomaly {
		voltage = 0
		motorTemp = 0
		motorRevolutions = 0
		gasPercentage = 10
		airFlow = 0
		airTemperature = 90
	}

	var event Event
	if eventType == AC_VENT {
		event = NewACEvent(airFlow, airTemperature, coolantTemperature)
		jsonBytes, _ := json.Marshal(event)
		return string(jsonBytes)

	} else if eventType == GENERATOR_EVENT {
		event = NewGeneratorEvent(hertz, amps, voltage, gasPercentage)
		jsonBytes, _ := json.Marshal(event)
		return string(jsonBytes)

	} else {
		event = NewMotorEvent(motorTemp, motorRevolutions)
		jsonBytes, _ := json.Marshal(event)
		return string(jsonBytes)
	}
}
```

### Azure SQL

#### Table Definitions

```sql
create table ACEvents
(
  Id int not null primary key identity,
  DeviceId varchar(10) not null,
  Ts DateTime not null,
  CoolantTemperature float not null,
  AirFlow float not null,
  AirTemperature float not null
)

create table GeneratorEvents
(
  Id int not null primary key identity,
  DeviceId varchar(10) not null,
  Ts DateTime not null,
  Hertz float not null,
  Amps float not null,
  Voltage float not null,
  GasPercentage float not null
)

create table MotorEvents
(
  Id int not null primary key identity,
  DeviceId varchar(10) not null,
  Ts DateTime not null,
  Temperature float not null,
  Revolutions float not null
)
```

### Stream Analytics Setup

#### Event Hubs Input

- Azure Hubs [hub-ecloud1-location1]
  - Hub [hub-location1]
    - Consumer group [hub_location1_cg]

#### SQL Output

Azure SQL Tables:

- ACMessages
- GeneratorMessages
- MotorMessages

#### Stream Analytics Jobs

> Note: One Stream Analytic jobs instance can process many jobs. The Stream Analytics query language can perform time based operations, aggregations, etc.

```
select a.deviceId,a.ts,a.coolantTemperature,a.airFlow,a.airTemperature
  into [hubdb-ACEvents] from [hub-ecloud1-location1] a 
  where type='ACEvent'

select a.deviceId,a.ts,a.hertz,a.amps,a.voltage,a.gasPercentage
  into [hubdb-GeneratorEvents] from [hub-ecloud1-location1] a 
  where type='GeneratorEvent'
  
select a.deviceId,a.ts,a.temperature,a.revolutions
  into [hubdb-MotorEvents] from [hub-ecloud1-location1] a 
  where type='MotorEvent'
```` 



