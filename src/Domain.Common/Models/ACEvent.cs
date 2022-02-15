namespace Domain.Common.Models
{
    public class ACMessage : BaseEvent
    {
        public ACEvent()
        {
            Type = typeof(ACEvent).Name;
        }
        public float CoolantTemperature { get; set; }
        public float AirFlow { get; set; }
        public float AirTemperature { get; set; }
    }
}
