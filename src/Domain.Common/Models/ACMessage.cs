namespace Domain.Common.Models
{
    public class ACMessage : BaseMessage
    {
        public ACMessage()
        {
            Type = typeof(ACMessage).Name;
        }
        public float CoolantTemperature { get; set; }
        public float AirFlow { get; set; }
        public float AirTemperature { get; set; }
    }
}
