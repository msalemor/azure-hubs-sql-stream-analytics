using System;

namespace Domain.Common.Models
{
    public class MotorEvent : BaseEvent
    {
        public MotorEvent()
        {
            Type = typeof(MotorEvent).Name;            
        }
        public float Temperature { get; set; }
        public float Revolutions { get; set; }
    }
}
