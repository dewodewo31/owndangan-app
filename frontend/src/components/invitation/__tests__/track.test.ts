import { track } from "../sections"

describe("track", () => {
  const url = "http://localhost:8080/api/v1/analytics/events"

  afterEach(() => {
    jest.restoreAllMocks()
  })

  it("emits the correct URL and JSON body via sendBeacon", async () => {
    const sendBeacon = jest.fn().mockReturnValue(true)
    Object.defineProperty(navigator, "sendBeacon", { value: sendBeacon, configurable: true })

    track("whatsapp_click", "evt-123")

    expect(sendBeacon).toHaveBeenCalledTimes(1)
    const [beaconUrl, blob] = sendBeacon.mock.calls[0]
    expect(beaconUrl).toBe(url)
    expect(blob).toBeInstanceOf(Blob)
    expect(blob.type).toBe("application/json")
    const blobText = await new Promise<string>((resolve) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result as string)
      reader.readAsText(blob)
    })
    expect(JSON.parse(blobText)).toEqual({ event_id: "evt-123", type: "whatsapp_click" })
  })

  it("falls back to fetch keepalive when sendBeacon is missing", () => {
    Object.defineProperty(navigator, "sendBeacon", { value: undefined, configurable: true })
    const fetchMock = jest.fn().mockResolvedValue({} as Response)
    global.fetch = fetchMock as unknown as typeof fetch

    track("map_click", "evt-123")

    expect(fetchMock).toHaveBeenCalledTimes(1)
    const [fetchUrl, init] = fetchMock.mock.calls[0]
    expect(fetchUrl).toBe(url)
    expect(init.method).toBe("POST")
    expect(init.keepalive).toBe(true)
    expect(JSON.parse(init.body)).toEqual({ event_id: "evt-123", type: "map_click" })
  })

  it("falls back to fetch when sendBeacon throws", () => {
    Object.defineProperty(navigator, "sendBeacon", {
      value: jest.fn().mockImplementation(() => {
        throw new Error("boom")
      }),
      configurable: true,
    })
    const fetchMock = jest.fn().mockResolvedValue({} as Response)
    global.fetch = fetchMock as unknown as typeof fetch

    expect(() => track("phone_click", "evt-123")).not.toThrow()
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })

  it("does nothing when eventId is empty", () => {
    const sendBeacon = jest.fn()
    Object.defineProperty(navigator, "sendBeacon", { value: sendBeacon, configurable: true })
    const fetchMock = jest.fn()
    global.fetch = fetchMock as unknown as typeof fetch

    track("map_click", undefined)
    track("map_click", "")

    expect(sendBeacon).not.toHaveBeenCalled()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it("never throws when sendBeacon and fetch both fail", () => {
    Object.defineProperty(navigator, "sendBeacon", {
      value: jest.fn().mockImplementation(() => {
        throw new Error("boom")
      }),
      configurable: true,
    })
    const fetchMock = jest.fn().mockRejectedValue(new Error("network down"))
    global.fetch = fetchMock as unknown as typeof fetch

    expect(() => track("whatsapp_click", "evt-123")).not.toThrow()
  })
})