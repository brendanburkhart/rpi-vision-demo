#include "pipeline.hpp"

Pipeline::Pipeline(int low, int high) : low_(low), high_(high) { }

void Pipeline::updateThresholds (int low, int high) {
    std::unique_lock<std::mutex> lock (mutex);
    low_ = low;
    high_ = high;
}

int Pipeline::low () {
    std::unique_lock<std::mutex> lock (mutex);
    return low_;
}

int Pipeline::high () {
    std::unique_lock<std::mutex> lock (mutex);
    return high_;
}
