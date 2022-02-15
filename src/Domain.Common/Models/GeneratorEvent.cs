namespace Domain.Common.Models
{

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
}
