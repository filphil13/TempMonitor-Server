import React, { useState, useEffect } from 'react';
import Graph from '../Graph/Graph';

const API_URL = "https://walrus-app-zu4le.ondigitalocean.app";
const PORT = "8080";
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
        fetch(API_URL + "/api/recent:" + PORT)
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
        <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
            <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
                <thead className="text-xs text-gray-700 uppercase dark:text-gray-400">
                    <tr>
                        <th scope="col" className="px-6 py-3 bg-gray-50 dark:bg-gray-800">
                            Name
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Temperature
                        </th>
                        <th scope="col" className="px-6 py-3 bg-gray-50 dark:bg-gray-800">
                            Humidity
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Time
                        </th>
                        <th scope="col" className="px-6 py-3">
                            Graph
                        </th>
                    </tr>
                </thead>
                <tbody>
                    {sensorList.length < 1 ? (
                        <tr className="border-b border-gray-200 dark:border-gray-700">
                            <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap bg-gray-50 dark:text-white dark:bg-gray-800">
                                No sensors found
                            </th>
                        </tr>
                    ) : (
                        sensorList.map((sensor) => (
                            <tr key={sensor.Name} className="border-b border-gray-200 dark:border-gray-700">
                                <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap bg-gray-50 dark:text-white dark:bg-gray-800">
                                    {sensor.Name}
                                </th>
                                <td className="px-6 py-4">
                                    {String(sensor.Temperature)}
                                </td>
                                <td className="px-6 py-4 bg-gray-50 dark:bg-gray-800">
                                    {String(sensor.Humidity)}
                                </td>
                                <td className="px-6 py-4">
                                    {String(Math.floor(Date.now() / 1000) - sensor.Time)} seconds ago
                                </td>
                            </tr>
                        ))
                    )}
                </tbody>
            </table>
        </div>
    );
}

export default TempTable;