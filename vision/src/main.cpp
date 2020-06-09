#include <iostream>
#include <thread>

#include <opencv2/opencv.hpp>
#include <cameraserver/CameraServer.h>

#include "pipeline.hpp"
#include "pipeline_controller.hpp"

int main ()
{
	Pipeline pipeline = Pipeline(10, 100);

	cs::CvSource outputStream = frc::CameraServer::GetInstance ()->PutVideo ("Gray", 640, 480);
	auto server = frc::CameraServer::GetInstance ()->StartAutomaticCapture (outputStream);
	std::cout << server.GetListenAddress() << ":" << server.GetPort() << std::endl;

	cv::VideoCapture cap;
	cap.open (0);

	if (!cap.isOpened ())
	{
		std::cerr << "Couldn't open capture." << std::endl;
		return -1;
	}

	cv::Mat frame;

	std::thread thread (&ServerImpl::Run, &pipeline);

	for (;;)
	{
		cap >> frame;
		if (frame.empty ()) break;

		cv::cvtColor (frame, frame, cv::COLOR_BGR2GRAY);
		cv::Canny (frame, frame, pipeline.low(), pipeline.high(), 3, true);

		outputStream.PutFrame(frame);
	}

	cap.release ();
	return 0;
}
