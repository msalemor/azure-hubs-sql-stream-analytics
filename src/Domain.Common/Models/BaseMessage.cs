using System;

namespace Domain.Common.Models
{

    //create table GeneratorMessages
    //(
    //	Id int not null primary key identity,
    //  DeviceId varchar(10) not null,
    //	Ts DateTime not null,
    //	Hertz float not null,
    //	Amps float not null,
    //	Voltage float not null,
    //	GasPercentage float not null
    //)

    //create table MotorMessages
    //(
    //	Id int not null primary key identity,
    //  DeviceId varchar(10) not null,
    //	Ts DateTime not null,
    //	Temperature float not null,
    //	Revolutions float not null
    //)

    //create table ACMessages
    //(
    //	Id int not null primary key identity,
    //  DeviceId varchar(10) not null,
    //	Ts DateTime not null,
    //	CoolantTemperature float not null,
    //	AirFlow float not null,
    //	AirTemperature float not null
    //)

    //SELECT
    //   a.deviceId DeviceId, a.ts Ts, a.temperature Temperature, a.revolutions Revolutions
    //INTO
    //  [ecloudhbdb - motormessages]
    //FROM
    
    //    [hubns-ecloud] a
    //WHERE type= 'MotorMessage'
    

    //SELECT
    //   a.deviceId DeviceId, a.ts Ts, a.airFlow AirFlow, a.airTemperature AirTemperature
    //INTO
    //  [ecloudhbdb - acmessages]
    //FROM
    
    //    [hubns-ecloud] a
    //WHERE type= 'ACMessage'
    

    //SELECT
    //    a.deviceId DeviceId, a.ts Ts, a.hertz Hertz, a.amps Amps, a.voltage Voltage, a.gasPercentage GasPercentage
    //INTO
    //  [ecloudhbdb - generatormessages]
    //FROM
    
    //    [hubns-ecloud] a
    //WHERE type= 'GeneratorMessage'

    public class BaseMessage
    {
        public BaseMessage()
        {
            Ts = DateTime.UtcNow;
        }

        public string DeviceId { get; set; }
        public string Type { get; set; }
        public DateTime Ts { get; set; }
    }
}
