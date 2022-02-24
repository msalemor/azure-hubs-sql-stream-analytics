with [allData] as (
    select * FROM [hub-ecloud1-location1]
),
anomalies AS (  
    SELECT a.ts,a.deviceId,'ACAnomality' as eventType, 'airFlow' as property, a.airflow as value 
    FROM allData a where a.type='ACEvent' and a.airFlow=0
    UNION
    SELECT a.ts,a.deviceId,'GeneratorAnomality' as eventType, 'voltage' as property, a.voltage as value 
    FROM allData a where a.type='GeneratorEvent' and a.voltage=0
    UNION
    SELECT a.ts,a.deviceId,'GeneratorAnomality' as eventType, 'gasPercentage' as property, a.voltage as value 
    FROM allData a where a.type='GeneratorEvent' and a.gasPercentage<20
    UNION
    SELECT a.ts,a.deviceId,'MotorAnomality' as eventType, 'revolutions' as property, a.revolutions as value 
    FROM allData a where a.type='MotorEvent' and a.revolutions=0
)

select 
  a.ts,a.deviceId,a.eventType,a.property,a.value
into
  [anomalies-hub]
from anomalies a;

select a.deviceId,a.ts,a.coolantTemperature,a.airFlow,a.airTemperature
  into [hubdb-ACEvents] from allData a 
  where type='ACEvent';

select a.deviceId,a.ts,a.hertz,a.amps,a.voltage,a.gasPercentage
  into [hubdb-GeneratorEvents] from allData a 
  where type='GeneratorEvent';
  
select a.deviceId,a.ts,a.temperature,a.revolutions
  into [hubdb-MotorEvents] from allData a 
  where type='MotorEvent';
