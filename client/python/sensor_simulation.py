#### SENSOR SIMULATOR ####
# @param: hour (0-4):
#   - 0: 4.01–9.00
#   - 1: 9.01–12.00
#   - 2: 12.01–17.00
#   - 3: 17.01-22.00
#   - 4: 22.01-4.00

temperature_matrix = [[4, 6, 8.0, 10, 4], 
                      [4, 6, 8.5, 11, 5], 
                      [5, 7, 10.0, 13, 6], 
                      [11, 13, 14.5, 16, 10], 
                      [15, 17, 18.5, 20, 15], 
                      [18, 20, 22.5, 25, 18], 
                      [24, 26, 28.0, 30, 23], 
                      [25, 27, 29.0, 31, 24], 
                      [17, 19, 21.5, 24, 17], 
                      [14, 16, 18.0, 20, 14], 
                      [13, 15, 17.0, 19, 10], 
                      [9, 11, 12.5, 14, 8]]

humidity_vector = [76, 76, 74, 77, 78, 76, 74, 75, 76, 77, 77, 78]


brightness_matrix = [[2000, 10000, 25000, 10000, 500],
                     [3000, 12000, 30000, 12000, 1000],  
                     [5000, 20000, 50000, 20000, 2000],   
                     [10000, 40000, 80000, 40000, 5000],  
                     [15000, 60000, 100000, 60000, 10000],
                     [20000, 80000, 120000, 80000, 15000],
                     [20000, 80000, 120000, 80000, 15000],
                     [18000, 70000, 110000, 70000, 12000],
                     [12000, 50000, 90000, 50000, 8000], 
                     [6000, 30000, 60000, 30000, 4000],   
                     [3000, 15000, 35000, 15000, 1000],
                     [2000, 10000, 25000, 10000, 500]]

pm10_matrix = [[65, 55, 45, 50, 60], 
               [60, 50, 40, 45, 55], 
               [55, 45, 35, 40, 50],
               [40, 30, 25, 30, 35],   
               [35, 25, 20, 25, 30],  
               [30, 20, 15, 20, 25],  
               [30, 20, 15, 20, 25],  
               [35, 25, 20, 25, 30],  
               [40, 30, 25, 30, 35],  
               [50, 40, 30, 35, 45],  
               [60, 50, 40, 45, 55],  
               [65, 55, 45, 50, 60]]

#@param month: The month of the year (0-11)
#@param hour: The hour of the day (0-4) 
def getTemperature(month, hour):
    """
    Function to get the temperature for a given month and hour.
    The temperature is calculated based on a predefined matrix.
    """

    return float(temperature_matrix[month][hour])

#@param month: The month of the year (0-11)
def getHumidity(month):
    """
    Function to get the humidity for a given month and hour.
    The humidity is calculated based on a predefined matrix.
    """
    return float(humidity_vector[month])

#@param month: The month of the year (0-11)
#@param hour: The hour of the day (0-4) 
def getBrightness(month, hour):
    """
    Function to get the brightness for a given hour.
    The brightness is calculated based on a predefined matrix.
    """
    return float(brightness_matrix[month][hour])


#@param month: The month of the year (0-11)
#@param hour: The hour of the day (0-4) 
def getAirQuality(month, hour):
    """
    Function to get the air quality for a given month and hour.
    The air quality is calculated based on a predefined matrix.
    """
    return float(pm10_matrix[month][hour])