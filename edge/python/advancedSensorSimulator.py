import math
import random
import datetime
import datetime
from zoneinfo import ZoneInfo

# Parametri globali per controllo rumore e anomalie
NOISE_FACTOR = 0.1
ANOMALY_PROBABILITY = 0.02
SENSOR_FAILURE_PROBABILITY = 0.01
ITALY_TIMEZONE = ZoneInfo("Europe/Rome")


def _parse_timestamp(timestamp_str):
    """Converte stringa timestamp in oggetto datetime e componenti temporali"""
    try:
        # Parsing del timestamp stringa
        dt = datetime.datetime.fromisoformat(timestamp_str.replace(' ', 'T'))
    except:
        # Fallback in caso di formato diverso
        try:
            dt = datetime.datetime.strptime(timestamp_str, '%Y-%m-%d %H:%M:%S.%f')
        except:
            dt = datetime.datetime.strptime(timestamp_str, '%Y-%m-%d %H:%M:%S')
    
    # Calcolo componenti temporali
    day_of_year = dt.timetuple().tm_yday
    hour_of_day = dt.hour + dt.minute / 60.0 + dt.second / 3600.0
    year_fraction = day_of_year / 365.25
    day_fraction = hour_of_day / 24.0
    
    return dt, year_fraction, day_fraction


def _add_noise_and_anomalies(base_value, noise_std, min_val=None, max_val=None):
    """Aggiunge rumore gaussiano e possibili anomalie al valore base"""
    # Rumore gaussiano
    noise = random.gauss(0, noise_std * NOISE_FACTOR)
    value = base_value + noise
    
    # Possibili anomalie (spike o drop)
    if random.random() < ANOMALY_PROBABILITY:
        anomaly_factor = random.choice([-1, 1]) * random.uniform(0.2, 0.5)
        value *= (1 + anomaly_factor)
    
    # Applica limiti se specificati
    if min_val is not None:
        value = max(value, min_val)
    if max_val is not None:
        value = min(value, max_val)
        
    return round(value, 1)


def get_temperature(timestamp_str):
    """
    Genera temperatura realistica basata su timestamp stringa
    
    Args:
        timestamp_str (str): Timestamp in formato string (es: '2024-12-20 14:30:25')
    
    Returns:
        float: Temperatura in gradi Celsius (o None se sensore guasto)
    
    Considera:
    - Ciclo stagionale annuale (freddo inverno, caldo estate)
    - Ciclo giornaliero (minimo all'alba, massimo pomeriggio)
    - Variazioni casuali e anomalie
    """
    dt, year_fraction, day_fraction = _parse_timestamp(timestamp_str)
    
    # Temperatura base stagionale (sinusoide annuale)
    # Minimo a gennaio (15 giorni), massimo a luglio (195 giorni)
    seasonal_temp = 15 + 12 * math.sin(2 * math.pi * (year_fraction - 0.04))
    
    # Variazione giornaliera (sinusoide giornaliera)
    # Minimo alle 6:00, massimo alle 14:00
    daily_variation = 8 * math.sin(2 * math.pi * (day_fraction - 0.25))
    
    # Temperatura base
    base_temp = seasonal_temp + daily_variation
    
    # Aggiunge rumore e anomalie
    temperature = _add_noise_and_anomalies(base_temp, 2.0, -10, 45)
    
    # Possibilità di valore None (sensore guasto)
    return None if random.random() < SENSOR_FAILURE_PROBABILITY else temperature


def get_humidity(timestamp_str):
    """
    Genera umidità realistica basata su timestamp stringa
    
    Args:
        timestamp_str (str): Timestamp in formato string
    
    Returns:
        float: Umidità relativa in percentuale (o None se sensore guasto)
    
    Considera:
    - Stagionalità (più alta in inverno, più bassa in estate)
    - Variazione giornaliera (più alta di notte)
    - Correlazione inversa con temperatura
    """
    dt, year_fraction, day_fraction = _parse_timestamp(timestamp_str)
    
    # Umidità base stagionale (inversa alla temperatura)
    seasonal_humidity = 75 - 15 * math.sin(2 * math.pi * (year_fraction - 0.04))
    
    # Variazione giornaliera (più alta di notte)
    daily_variation = 10 * math.sin(2 * math.pi * (day_fraction + 0.5))
    
    # Umidità base
    base_humidity = seasonal_humidity + daily_variation
    
    # Correlazione con temperatura (effetto evaporazione)
    temp = get_temperature(timestamp_str)
    if temp is not None and temp > 25:
        base_humidity -= (temp - 25) * 0.5  # Diminuisce con alta temperatura
    
    # Aggiunge rumore e anomalie
    humidity = _add_noise_and_anomalies(base_humidity, 5.0, 20, 95)
    
    # Possibilità di valore None
    return None if random.random() < SENSOR_FAILURE_PROBABILITY else humidity


def get_brightness(timestamp_str):
    """
    Genera luminosità realistica basata su timestamp stringa
    
    Args:
        timestamp_str (str): Timestamp in formato string
    
    Returns:
        float: Luminosità in lux (o None se sensore guasto)
    
    Considera:
    - Ciclo giornaliero (alba, mezzogiorno, tramonto)
    - Variazione stagionale (giornate più lunghe in estate)
    - Condizioni meteorologiche simulate (nuvolosità)
    """
    dt, year_fraction, day_fraction = _parse_timestamp(timestamp_str)
    
    # Durata del giorno stagionale (più lungo in estate)
    day_length_factor = 0.8 + 0.4 * math.sin(2 * math.pi * (year_fraction - 0.04))
    
    # Curva di luminosità giornaliera (campana gaussiana)
    # Centrata a mezzogiorno (0.5)
    hour_from_noon = abs(day_fraction - 0.5)
    
    if hour_from_noon < day_length_factor * 0.3:  # Ore di luce
        # Luminosità massima a mezzogiorno
        brightness_factor = math.exp(-((hour_from_noon * 8) ** 2))
        base_brightness = 100000 * brightness_factor
    else:
        # Notte o crepuscolo
        base_brightness = 50 + 500 * math.exp(-((hour_from_noon - 0.3) * 20) ** 2)
    
    # Variazione stagionale (sole più alto in estate)
    seasonal_factor = 0.7 + 0.3 * math.sin(2 * math.pi * (year_fraction - 0.04))
    base_brightness *= seasonal_factor
    
    # Simulazione nuvolosità (riduce luminosità)
    if random.random() < 0.3:  # 30% probabilità di nuvolosità
        cloud_factor = random.uniform(0.2, 0.8)
        base_brightness *= cloud_factor
    
    # Aggiunge rumore e anomalie
    brightness = _add_noise_and_anomalies(base_brightness, base_brightness * 0.1, 0, 150000)
    
    # Possibilità di valore None
    return None if random.random() < SENSOR_FAILURE_PROBABILITY else brightness


def get_air_quality(timestamp_str):
    """
    Genera qualità dell'aria (PM10) realistica basata su timestamp stringa
    
    Args:
        timestamp_str (str): Timestamp in formato string
    
    Returns:
        float: Concentrazione PM10 in μg/m³ (o None se sensore guasto)
    
    Considera:
    - Stagionalità (peggiore in inverno per riscaldamento)
    - Variazione giornaliera (peggiore nelle ore di traffico)
    - Condizioni meteorologiche (vento, pioggia)
    """
    dt, year_fraction, day_fraction = _parse_timestamp(timestamp_str)
    
    # PM10 base stagionale (più alto in inverno)
    seasonal_pm10 = 35 + 20 * math.sin(2 * math.pi * (year_fraction + 0.5))
    
    # Variazione giornaliera (picchi nelle ore di traffico)
    # Picco mattutino (8:00) e serale (18:00)
    morning_peak = 15 * math.exp(-((day_fraction - 0.33) * 12) ** 2)  # 8:00
    evening_peak = 20 * math.exp(-((day_fraction - 0.75) * 12) ** 2)   # 18:00
    traffic_variation = morning_peak + evening_peak
    
    # PM10 base
    base_pm10 = seasonal_pm10 + traffic_variation
    
    # Effetto meteorologico simulato
    weather_factor = 1.0
    if random.random() < 0.2:  # 20% probabilità di pioggia (riduce PM10)
        weather_factor = random.uniform(0.3, 0.7)
    elif random.random() < 0.1:  # 10% probabilità di alta pressione (aumenta PM10)
        weather_factor = random.uniform(1.3, 2.0)
    
    base_pm10 *= weather_factor
    
    # Aggiunge rumore e anomalie
    pm10 = _add_noise_and_anomalies(base_pm10, 8.0, 5, 200)
    
    # Possibilità di valore None
    return None if random.random() < SENSOR_FAILURE_PROBABILITY else pm10

def get_timestamp():
    now_in_italy = datetime.datetime.now(ITALY_TIMEZONE)
    timestamp = now_in_italy.strftime("%Y-%m-%d %H:%M:%S.%f")
    return timestamp
