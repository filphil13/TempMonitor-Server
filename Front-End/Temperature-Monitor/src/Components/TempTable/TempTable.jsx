import React, { useState, useEffect } from 'react';

const API_URL = "https://temp-monitor-a38f32c02c5e.herokuapp.com";

function TempTable() {
    const [sensorList, setSensorList] = useState([]);

    useEffect(() => {
        // Fetch recent data every 3 seconds
        const interval = setInterval(() => {
            getRecentData();
        }, 3000);

        // Clean up the interval on component unmount
        return () => clearInterval(interval);
    }, []);

    function getRecentData() {
        fetch(API_URL + "/api/recent")
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
            });
    }

    return (
        <div className="overflow-x-auto sm:overflow-visible">
            <div className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
                <div className="text-xs text-gray-700 uppercase dark:text-gray-400">
                    <div className="flex flex-wrap border-b border-gray-200 dark:border-gray-700">
                        <div className="w-full px-6 py-3 bg-gray-50 dark:bg-gray-800 sm:w-auto sm:border-r md:w-1/4">
                            Name
                        </div>
                        <div className="w-full px-6 py-3 sm:w-auto sm:border-r md:w-1/4">
                            Temperature
                        </div>
                        <div className="w-full px-6 py-3 bg-gray-50 dark:bg-gray-800 sm:w-auto sm:border-r md:w-1/4">
                            Humidity
                        </div>
                        <div className="w-full px-6 py-3 sm:w-auto md:w-1/4">
                            Time
                        </div>
                    </div>
                </div>
                {sensorList.length < 1 ? (
                    <div className="flex flex-wrap border-b border-gray-200 dark:border-gray-700">
                        <div className="w-full px-6 py-4 font-medium text-gray-900 whitespace-nowrap bg-gray-50 dark:text-white dark:bg-gray-800">
                            No sensors found
                        </div>
                    </div>
                ) : (
                    sensorList.map((sensor) => (
                        <div key={sensor.Name} className="flex flex-wrap border-b border-gray-200 dark:border-gray-700">
                            <div className="w-full px-6 py-4 font-medium text-gray-900 whitespace-nowrap bg-gray-50 dark:text-white dark:bg-gray-800 sm:w-auto sm:border-r md:w-1/4">
                                {sensor.Name}
                            </div>
                            <div className="w-full px-6 py-4 sm:w-auto sm:border-r md:w-1/4">
                                {String(sensor.Temperature)}
                            </div>
                            <div className="w-full px-6 py-4 bg-gray-50 dark:bg-gray-800 sm:w-auto sm:border-r md:w-1/4">
                                {String(sensor.Humidity)}
                            </div>
                            <div className="w-full px-6 py-4 sm:w-auto md:w-1/4">
                                {sensor.Time}
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}

export default TempTable;