using System;

namespace Domain.Common.Models
{
    public class MotorMessage : BaseMessage
    {
        public MotorMessage()
        {
            Type = typeof(MotorMessage).Name;            
        }
        public float Temperature { get; set; }
        public float Revolutions { get; set; }
    }
}
