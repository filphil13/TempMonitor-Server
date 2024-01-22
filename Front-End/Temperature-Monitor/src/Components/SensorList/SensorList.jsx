import React, { useState, useEffect } from 'react';
import SensorCard from './SensorCard';

const API_URL = "https://oyster-app-rwyik.ondigitalocean.app";
const PORT = "443";
function SensorList() {
    const [sensorList, setSensorList] = useState([]);
    
    useEffect(() => {
        // Fetch recent data every 3 seconds
        const interval = setInterval(() => {
            getSensorData();
        }, 3000);

        // Clean up the interval on component unmount
        return () => clearInterval(interval);
    }, []);


    function getSensorData() {
        fetch(API_URL + ":" + PORT + "/api/sensors?userToken=1234567890")
            .then(async response => {
                const data = await response.json();

                // Check for error response
                if (!response.ok) {
                    // Get error message from body or default to response statusText
                    const error = (data && data.message) || response.statusText;
                    return Promise.reject(error);
                } else if (data.length === 0) {
                    return;
                }

                setSensorList(data);
            })
            .catch(error => console.error('Fetch error:', error));;
    }

    return (
        <>
            {sensorList.length < 1 ? (
           <p> No sensors found </p>
            ) : (
            <ul className='list-none w-screen'>
                {sensorList.map((sensor) => (
                        <li key={sensor.Name} className=" min-w-screen bg-black shadow-md border border-gray-200 rounded-lg dark:bg-gray-800 dark:border-gray-700" >
                            <SensorCard  
                                name={sensor.Name} 
                                temperature={sensor.Temperature} 
                                humidity={sensor.Humidity}
                                status={sensor.Status}
                            />
                            </li>
                ))}
            </ul>
            )}
        </>                    
    );
}

export default SensorList;