import React, { useEffect, useRef } from 'react';
import {Chart} from 'chart.js';

const Graph = ({ sensorData }) => {
    const chartRef = useRef(null);

    useEffect(() => {
        const ctx = chartRef.current.getContext('2d');
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: sensorData.Log.map(data => data.date),
                datasets: [
                    {
                        type: 'line',
                        label: 'Temperature',
                        data: sensorData.map(data => data.Temperature),
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1,
                    },
                    {
                        type: 'line',
                        label: 'Humidity',
                        data: sensorData.map(data => data.Humidity),
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
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
    }, [sensorData]);

    return <canvas ref={chartRef} />;
};

export default Graph;
