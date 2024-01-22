import React, { useEffect, useRef, useState } from 'react';
import { Chart, LinearScale } from 'chart.js'; // Import the LinearScale module

const Graph = ({ sensorName }) => {
    const chartRef = useRef(null);
    const [tempLog, setTempLog] = useState([]);

    useEffect(() => {
        const interval = setInterval(() => {
            getAllSensorLogs();
        }, 3000);

        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        if (sensorName) {
            getAllSensorLogs();
        }
    }, [sensorName]);


    const getAllSensorLogs = () => {
        fetch(API_URL + "/api/all?" + sensorName)
            .then(async (response) => {
                const data = await response.json();
                if (!response.ok) {
                    const error = (data && data.message) || response.statusText;
                    return Promise.reject(error);
                } else if (data.length === 0) {
                    return;
                }
                setTempLog(data);
            })
            .catch((error) => {
                console.error("Error fetching sensor data:", error);
            });
    };

    useEffect(() => {
        const ctx = chartRef.current.getContext('2d');

        new Chart(ctx, {
            type: 'line',
            data: {
                labels: tempLog.map((data) => data.Time),
                datasets: [
                    {
                        type: 'line',
                        label: 'Temperature',
                        data: tempLog.map((data) => data.Temperature),
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1,
                    },
                    {
                        type: 'line',
                        label: 'Humidity',
                        data: tempLog.map((data) => data.Humidity),
                        backgroundColor: 'rgba(192, 75, 75, 0.2)',
                        borderColor: 'rgba(192, 75, 75, 1)',
                        borderWidth: 1,
                    },
                ],
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true,
                    },
                },
            },
        });
    }, [tempLog]);

    return <canvas ref={chartRef} />;
};

export default Graph;
