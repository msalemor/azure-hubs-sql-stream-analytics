namespace Domain.Common.Models
{

    public class GeneratorMessage : BaseMessage
    {
        public GeneratorMessage()
        {
            Type = typeof(GeneratorMessage).Name;
        }
        public float Hertz { get; set; }
        public float Amps { get; set; }
        public float Voltage { get; set; }
        public float GasPercentage { get; set; }

    }
}
