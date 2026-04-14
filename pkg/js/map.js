function getFlowColor(value) {
  var absValue = Math.abs(value);
  var ratio = Math.min(1, Math.max(0, absValue));

  var g = Math.round(255 * (1 - ratio));

  return "rgb(255, " + g + ", 0)";
}

function updateFlowValues(flowData, lineByMrid, map, flowLayer) {
  flowLayer.clearLayers();

  for (var mrid in flowData) {
    var value = flowData[mrid];
    if (!value) continue;

    var line = lineByMrid[mrid];
    if (!line) continue;

    var midLat = (line.LatFrom + line.LatTo) / 2;
    var midLng = (line.LngFrom + line.LngTo) / 2;

    var isPositive = value >= 0;
    var deltaLat = line.LatTo - line.LatFrom;
    var deltaLng = line.LngTo - line.LngFrom;
    var isVertical = Math.abs(deltaLat) > Math.abs(deltaLng);

    var arrowSymbol;
    if (isVertical) {
      arrowSymbol = isPositive
        ? deltaLat > 0
          ? "↓"
          : "↑"
        : deltaLat > 0
          ? "↑"
          : "↓";
    } else {
      arrowSymbol = isPositive ? "→" : "←";
    }

    var displayValue = Math.abs(parseFloat(value)).toFixed(1);

    var html =
      '<span class="flow-arrow">' +
      arrowSymbol +
      '</span><span class="flow-value">' +
      displayValue +
      "</span>";

    var label = L.divIcon({
      className: "flow-label",
      html: html,
      iconSize: [50, 24],
      iconAnchor: [25, 12],
    });

    var marker = L.marker([midLat, midLng], { icon: label });
    flowLayer.addLayer(marker);
    marker.getElement().style.backgroundColor = color;
  }
}

function initMap(substations, lines) {
  var map = L.map("map").setView([59.9139, 10.7522], 6);

  map.on("popupopen", function (e) {
    var content = e.popup.getElement().querySelector(".leaflet-popup-content");
    if (content) {
      htmx.process(content);
    }
  });

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution:
      '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
  }).addTo(map);

  var voltageLevels = [
    { min: 0, max: 66, color: "#808080", label: "< 66 kV" },
    { min: 66, max: 132, color: "#4169E1", label: "66-132 kV" },
    { min: 132, max: 220, color: "#000000", label: "132-220 kV" },
    { min: 220, max: 300, color: "#9370DB", label: "220-300 kV" },
    { min: 300, max: 380, color: "#D5B60A", label: "220-300 kV" },
    { min: 380, max: Infinity, color: "#DC143C", label: "> 380 kV" },
  ];

  var voltageControls = document.getElementById("voltage-controls");
  var lineLayers = {};
  var lineVisibility = {};

  voltageLevels.forEach(function (level, index) {
    var group = document.createElement("div");
    group.className = "voltage-group";

    var label = document.createElement("label");
    var checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.checked = true;
    checkbox.id = "voltage-" + index;
    checkbox.addEventListener("change", function () {
      toggleVoltageLevel(index, this.checked);
    });

    var colorIndicator = document.createElement("div");
    colorIndicator.className = "color-indicator";
    colorIndicator.style.backgroundColor = level.color;

    var voltageLabel = document.createElement("span");
    voltageLabel.className = "voltage-label";
    voltageLabel.textContent = level.label;

    label.appendChild(checkbox);
    label.appendChild(colorIndicator);
    label.appendChild(voltageLabel);
    group.appendChild(label);
    voltageControls.appendChild(group);

    lineLayers[index] = L.layerGroup().addTo(map);
    lineVisibility[index] = true;
  });

  var flowLayers = {};
  voltageLevels.forEach(function (level, index) {
    flowLayers[index] = L.layerGroup().addTo(map);
  });

  var substationMarkers = [];

  substations.forEach(function (sub) {
    var circle = L.circleMarker([sub.Lat, sub.Lng], {
      radius: 5,
      fillColor: "#3498db",
      color: "#2980b9",
      weight: 2,
      opacity: 0.8,
      fillOpacity: 0.7,
    })
      .addTo(map)
      .bindPopup(
        "<p><strong>" +
          sub.Name +
          "</strong></p><p>" +
          sub.Mrid +
          '</p><button class="button is-small is-primary" hx-post="/production?mrid=' +
          sub.Mrid +
          "&name=" +
          encodeURIComponent(sub.Name) +
          '" hx-include="#active-production-form" hx-swap="innerHTML" hx-target="#active-production-form">Produce</button>',
      );

    substationMarkers.push(circle);
  });

  document.getElementById("substation-count").textContent = substations.length;

  var totalLines = lines.length;

  var linesByLevel = {};
  voltageLevels.forEach(function (level, index) {
    linesByLevel[index] = [];
  });

  lines.forEach(function (line) {
    var voltage = line.Voltage || 0;
    var levelIndex = getVoltageLevelIndex(voltage);
    linesByLevel[levelIndex].push(line);
  });

  for (var i = 0; i < voltageLevels.length; i++) {
    linesByLevel[i].forEach(function (line) {
      var latlngs = [
        [line.LatFrom, line.LngFrom],
        [line.LatTo, line.LngTo],
      ];

      var polyline = L.polyline(latlngs, {
        color: voltageLevels[i].color,
        weight: 3,
        opacity: 0.8,
      }).bindPopup(
        "<strong>" +
          line.Name +
          "</strong><br>" +
          line.Voltage +
          " kV<br>" +
          line.Mrid,
      );

      lineLayers[i].addLayer(polyline);
    });
  }

  var lineByMrid = {};
  lines.forEach(function (line) {
    lineByMrid[line.Mrid] = line;
  });

  document.addEventListener("action-form-changed", function () {
    var form = document.getElementById("active-production-form");
    var formData = new FormData(form);
    var urlEncoded = new URLSearchParams(formData).toString();

    fetch("/flow", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: urlEncoded,
    })
      .then(function (response) {
        return response.json();
      })
      .then(function (flowData) {
        updateFlowValues(flowData.flow, lineByMrid, map, flowLayers);
      });
  });

  document.getElementById("line-count").textContent = totalLines;

  function getVoltageLevelIndex(voltage) {
    for (var i = 0; i < voltageLevels.length; i++) {
      if (voltage >= voltageLevels[i].min && voltage < voltageLevels[i].max) {
        return i;
      }
    }
    return voltageLevels.length - 1;
  }

  function toggleVoltageLevel(levelIndex, visible) {
    lineVisibility[levelIndex] = visible;
    if (visible) {
      map.addLayer(lineLayers[levelIndex]);
    } else {
      map.removeLayer(lineLayers[levelIndex]);
    }
  }

  function voltageColor(voltage) {
    var levelIndex = getVoltageLevelIndex(voltage);
    return voltageLevels[levelIndex].color;
  }
}
