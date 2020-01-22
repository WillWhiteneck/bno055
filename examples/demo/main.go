package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/willwhiteneck/bno055"
)

func main() {
	sensor, err := bno055.NewSensor(0x28, 3)
	if err != nil {
		panic(err)
	}

	err = sensor.UseExternalCrystal(true)
	if err != nil {
		panic(err)
	}

	var (
		isCalibrated       bool
		calibrationOffsets bno055.CalibrationOffsets
		calibrationStatus  *bno055.CalibrationStatus
	)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for !isCalibrated {
		select {
		case <-signals:
			err := sensor.Close()
			if err != nil {
				panic(err)
			}
		default:
			calibrationOffsets, calibrationStatus, err = sensor.Calibration()
			if err != nil {
				panic(err)
			}

			isCalibrated = calibrationStatus.IsCalibrated()

			fmt.Printf(
				"\r*** Calibration status (0..3): system=%v, accelerometer=%v, gyroscope=%v, magnetometer=%v",
				calibrationStatus.System,
				calibrationStatus.Accelerometer,
				calibrationStatus.Gyroscope,
				calibrationStatus.Magnetometer,
			)
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("*** Done! Calibration offsets: %v\n", calibrationOffsets)

	status, err := sensor.Status()
	if err != nil {
		panic(err)
	}

	fmt.Printf("*** Status: system=%v, system_error=%v, self_test=%v\n", status.System, status.SystemError, status.SelfTest)

	revision, err := sensor.Revision()
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"*** Revision: software=%v, bootloader=%v, accelerometer=%v, gyroscope=%v, magnetometer=%v\n",
		revision.Software,
		revision.Bootloader,
		revision.Accelerometer,
		revision.Gyroscope,
		revision.Magnetometer,
	)

	// axisConfig, err := sensor.AxisConfig()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf(
	// 	"*** Axis: x=%v, y=%v, z=%v, sign_x=%v, sign_y=%v, sign_z=%v\n",
	// 	axisConfig.X,
	// 	axisConfig.Y,
	// 	axisConfig.Z,
	// 	axisConfig.SignX,
	// 	axisConfig.SignY,
	// 	axisConfig.SignZ,
	// )

	temperature, err := sensor.Temperature()
	if err != nil {
		panic(err)
	}

	fmt.Printf("*** Temperature: t=%v\n", temperature)

	for {
		select {
		case <-signals:
			err := sensor.Close()
			if err != nil {
				panic(err)
			}
		default:
			vector, err := sensor.Euler()
			if err != nil {
				panic(err)
			}

			fmt.Printf("\r*** Euler angles: x=%5.3f, y=%5.3f, z=%5.3f", vector.X, vector.Y, vector.Z)
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Output:
	// *** Status: system=133, system_error=0, self_test=15
	// *** Revision: software=785, bootloader=21, accelerometer=251, gyroscope=15, magnetometer=50
	// *** Axis: x=0, y=1, z=2, sign_x=0, sign_y=0, sign_z=0
	// *** Temperature: t=27
	// *** Euler angles: x=2.312, y=2.000, z=91.688
}
