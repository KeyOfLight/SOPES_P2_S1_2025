from locust import HttpUser, task, between
import random

# Crear un arreglo de 10,000 datos de ejemplo
weather_data_array = [
    {
        "Description": f"Descripción {i}",
        "Country": "GT",
        "Weather": random.choice(["Soleado", "Lluvioso", "Nublado", "Tormentoso"])
    }
    for i in range(10000)
]

class WeatherUser(HttpUser):
    wait_time = between(0.5, 1)  # Menor tiempo de espera para bombardear más rápido

    @task
    def send_weather_data(self):
        # Seleccionar un dato aleatorio del arreglo
        payload = random.choice(weather_data_array)
        self.client.post("/input", json=payload)


# locust -f locustfile.py --host http://127.0.0.1.nip.io
