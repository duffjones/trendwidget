//These functions control the settings, rendering and animation for each of our charts. 
//the data is stored at the file location called by 'CSVURL'

$(function() {
    var myChart = Highcharts.chart('thirtydaytemp', {
        chart: {
            type: 'area',
            backgroundColor: 'rgba(255, 255, 255, 0.0)',
            color: 'white'
        },
        title: {
            text: 'Thirty Day Temperature',
            style: {
                color: 'white'
            }
        },

        subtitle: {
            text: '',
            color: 'white'
        },

        plotOptions: {
            series: {
                color: 'white'
            }
        },

        xAxis: {
            labels: {
                style: {
                    color: 'white',
                }
            },
        },
        yAxis: {

            labels: {
                style: {
                    color: 'white',
                    font: '11px Trebuchet MS, Verdana, sans-serif'
                }
            },
        },
        data: {
            csvURL: 'http://localhost:8080/csv/CurrentTempCelcius_30days_trend.csv',
            enablePolling: true,
            color: 'white'
        }

    });
});

$(function() {
    var myChart2 = Highcharts.chart('fivedaytemp', {
        chart: {
            type: 'spline',
            backgroundColor: 'rgba(255, 255, 255, 0.0)'

        },
        title: {
            style: {
                color: 'white'
            },
            text: 'Temperature over 5 Days'
        },

        plotOptions: {
            series: {
                color: 'white'
            }
        },

        xAxis: {
            labels: {
                style: {
                    color: 'white',
                }
            },
        },
        yAxis: {

            labels: {
                style: {
                    color: 'white',
                    font: '11px Trebuchet MS, Verdana, sans-serif'
                }
            },
        },
        subtitle: {
            text: ''
        },

        data: {
            csvURL: 'http://localhost:8080/csv/CurrentTempCelcius_5days_trend.csv',
            enablePolling: true
        }
    });
});


$(function() {
    var myChart3 = Highcharts.chart('thirtydayhumid', {
        chart: {
            type: 'spline',
            backgroundColor: 'rgba(255, 255, 255, 0.0)',
            color: 'white'
        },
        title: {
            text: 'Thirty Day Humidity',
            style: {
                color: 'white'
            }
        },

        subtitle: {
            text: '',
            color: 'white'
        },

        plotOptions: {
            series: {
                color: 'white'
            }
        },

        xAxis: {
            labels: {
                style: {
                    color: 'white',
                }
            },
        },
        yAxis: {

            labels: {
                style: {
                    color: 'white',
                    font: '11px Trebuchet MS, Verdana, sans-serif'
                }
            },
        },
        data: {
            csvURL: 'http://localhost:8080/csv/Humidity_30days_trend.csv',
            enablePolling: true,
            color: 'white'
        }
    });
});




$(function() {
    var myChart4 = Highcharts.chart('thirtydaywind', {
        chart: {
            type: 'spline',
            backgroundColor: 'rgba(255, 255, 255, 0.0)',
            color: 'white'
        },
        title: {
            text: 'Thirty Day Windspeed',
            style: {
                color: 'white'
            }
        },

        subtitle: {
            text: '',
            color: 'white'
        },

        plotOptions: {
            series: {
                color: 'white'
            }
        },

        xAxis: {
            labels: {
                style: {
                    color: 'white',
                }
            },
        },
        yAxis: {
            labels: {
                style: {
                    color: 'white',
                    font: '11px Trebuchet MS, Verdana, sans-serif'
                }
            },
        },
        data: {
            csvURL: 'http://localhost:8080/csv/Windspeed_30days_trend.csv',
            enablePolling: true,
            color: 'white'
        }
    });


});


$(function() {
    var myChart5 = Highcharts.chart('fivedaywind', {
        chart: {
            backgroundColor: 'rgba(255, 255, 255, 0.0)',
            type: 'spline'
        },
        title: {
            text: '5 Day Windspeed',
            style: {
                color: 'white'
            },
        },
        plotOptions: {
            series: {
                color: 'white'
            }
        },

        xAxis: {
            title: false,
            labels: {
                style: {
                    color: 'white',
                }
            },
        },
        yAxis: {
            title: false,
            labels: {
                style: {
                    color: 'white',
                    font: '11px Trebuchet MS, Verdana, sans-serif'
                }
            },
        },
        subtitle: {
            text: ''
        },

        data: {
            csvURL: 'http://localhost:8080/csv/Windspeed_5days_trend.csv',
            enablePolling: true
        }
    });


});

$(function() {
    var myChart6 = Highcharts.chart('fivedayhumid', {
        chart: {
            type: 'spline',
            backgroundColor: 'rgba(255, 255, 255, 0.0)'

        },
        title: {
            text: 'Five Days Humidity',
            style: {
                color: 'white'
            }
        },
        plotOptions: {
            series: {
                color: 'white'
            }
        },

        xAxis: {
            title: false,
            labels: {
                style: {
                    color: 'white',
                }
            },
        },
        yAxis: {
            title: false,
            labels: {
                style: {
                    color: 'white',
                    font: '11px Trebuchet MS, Verdana, sans-serif'
                }
            },
        },

        subtitle: {
            text: ''
        },

        data: {
            csvURL: 'http://localhost:8080/csv/Humidity_5days_trend.csv',
            enablePolling: true
        }
    });
});