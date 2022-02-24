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
